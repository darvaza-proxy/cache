package memcache

import (
	"context"
	"sync"
	"time"

	"darvaza.org/cache"
	"darvaza.org/core"
	"darvaza.org/slog"
)

var (
	_ cache.Getter[string] = (*SingleFlight[string])(nil)
	_ cache.Setter[string] = (*SingleFlight[string])(nil)

	_ cache.Getter[uint64] = (*SingleFlight[uint64])(nil)
	_ cache.Setter[uint64] = (*SingleFlight[uint64])(nil)

	_ cache.Getter[[32]byte] = (*SingleFlight[[32]byte])(nil)
	_ cache.Setter[[32]byte] = (*SingleFlight[[32]byte])(nil)
)

// SingleFlight is a man-in-the-middle between [cache.Cache] and
// our [LRU] to prevent stampedes
type SingleFlight[K comparable] struct {
	name    string
	mu      sync.Mutex
	log     slog.Logger
	inward  AdderGetter[K]
	outward cache.Getter[K]
	getters map[K]*outreacher[K]
}

// NewSingleFlight creates a new [SingleFlight] controller, with an [LRU] for
// local cache and a [cache.Getter] to acquire the data externally.
// [SingleFlight] will prevent multiple requests for the same key to reach out
// at the same time.
func NewSingleFlight[K comparable](name string, inward AdderGetter[K], outward cache.Getter[K]) *SingleFlight[K] {
	if inward == nil || outward == nil {
		core.Panic("missing parameters")
	}

	sf := &SingleFlight[K]{
		name:    name,
		inward:  inward,
		outward: outward,
		getters: make(map[K]*outreacher[K]),
	}

	return sf
}

// Name returns the name of the Cache namespace
func (sf *SingleFlight[K]) Name() string {
	return sf.name
}

// revive:disable:cognitive-complexity
// revive:disable:cyclomatic
// revive:disable:function-length

// Get attempts to get the value of a key from its internal cache, otherwise reaches
// out to the provided [cache.Getter], but only once. While this is in process any other
// request for the same key will be held until we have a response from from the first.
func (sf *SingleFlight[K]) Get(ctx context.Context, key K, dest cache.Sink) error {
	// revive:enable:cognitive-complexity
	// revive:enable:cyclomatic
	// revive:enable:function-length
	sf.mu.Lock()

	if v, ex, ok := sf.inward.Get(key); ok {
		// cache hit
		defer sf.mu.Unlock()

		if log, ok := sf.withDebug(); ok {
			log.WithField("key", key).
				Print("hit")
		}

		var expire time.Time
		if ex != nil {
			expire = *ex
		}
		return dest.SetBytes(v, expire)
	}

	// cache miss
	if log, ok := sf.withDebug(); ok {
		log.WithField("key", key).
			Print("miss")
	}

	cond, first := sf.getCond(key)
	if !first {
		// wait patiently
		if log, ok := sf.withDebug(); ok {
			log.WithField("key", key).
				Print("waiting...")
		}
		cond.Wait()

		defer sf.mu.Unlock()
		if err := cond.Err(); err != nil {
			// failed
			if log, ok := sf.withDebug(); ok {
				log.WithField("key", key).
					Println("failed:", err)
			}
			return err
		}

		if log, ok := sf.withDebug(); ok {
			log.WithField("key", key).
				Print("ready")
		}

		// store copy on the [cache.Sink]
		return dest.SetBytes(cond.Bytes(), cond.Expire())
	}

	// reach out
	sf.mu.Unlock()
	if log, ok := sf.withDebug(); ok {
		log.WithField("key", key).
			Print("getting...")
	}
	err := sf.outward.Get(ctx, key, dest)
	// and lock again
	sf.mu.Lock()
	defer sf.mu.Unlock()

	if err == nil {
		// successfully acquired the value
		value, expire := dest.Bytes(), dest.Expire()

		if log, ok := sf.withDebug(); ok {
			log.WithField("key", key).
				WithField("size", len(value)).
				Print("success")
		}

		// store inward
		sf.inward.Add(key, value, expire)
		// and on the condition variable if anyone is waiting
		if !cond.Done() {
			cond.SetBytes(value, expire)
		}
		return nil
	}

	if cond.Ok() {
		// someone provided the value for us. happy days
		value, expire := cond.Bytes(), cond.Expire()

		if log, ok := sf.withDebug(); ok {
			log.WithField("key", key).
				WithField("size", len(value)).
				Print("thank you!")
		}

		return dest.SetBytes(value, expire)
	}

	// failed to acquire a value
	if log, ok := sf.withDebug(); ok {
		log.WithField("key", key).
			Println("failed:", err)
	}

	cond.SetError(err)
	return err
}

// Set stores the value for a key inward, and shares it with anyway waiting for it
func (sf *SingleFlight[K]) Set(_ context.Context, key K, value []byte,
	expire time.Time, _ cache.Type) error {
	//
	sf.mu.Lock()
	defer sf.mu.Unlock()

	sf.inward.Add(key, value, expire)
	if p, ok := sf.getters[key]; ok {
		defer p.Done()

		// there is people waiting
		p.SetBytes(value, expire)
	}
	return nil
}

// SetLogger attaches a [slog.Logger] to this [SingleFlight] quasi-[cache.Cache]
func (sf *SingleFlight[K]) SetLogger(log slog.Logger) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	sf.log = log
}

func (sf *SingleFlight[K]) withLogger(level slog.LogLevel) (slog.Logger, bool) {
	if sf.log != nil {
		log, ok := sf.log.WithLevel(level).WithEnabled()
		if ok {
			log = log.WithField("cache", sf.name)
			return log, true
		}
	}
	return nil, false
}

func (sf *SingleFlight[K]) withDebug() (slog.Logger, bool) {
	return sf.withLogger(slog.Debug)
}

func (sf *SingleFlight[K]) getCond(key K) (*outreacher[K], bool) {
	p, ok := sf.getters[key]
	if ok {
		// old
		return p, false
	}

	// new
	p = &outreacher[K]{
		parent: sf,
		cond:   sync.NewCond(&sf.mu),
		key:    key,
	}
	sf.getters[key] = p

	return p, true
}

type outreacher[K comparable] struct {
	parent *SingleFlight[K]
	count  int
	cond   *sync.Cond
	key    K

	done bool
	err  error
	b    []byte
	ex   time.Time
}

// Err returns the error set by a failed outward.Get()
func (p *outreacher[K]) Err() error { return p.err }

// Bytes returns the data set by a successful outward.Get()
func (p *outreacher[K]) Bytes() []byte { return p.b }

// Expire returns the expiration date set by a successful outward.Get()
func (p *outreacher[K]) Expire() time.Time { return p.ex }

// Ok tells if a value has been stored
func (p *outreacher[K]) Ok() bool {
	return p.done && p.err == nil
}

// SetBytes stores the result of a successful outward.Get()
func (p *outreacher[K]) SetBytes(v []byte, e time.Time) {
	defer p.cond.Broadcast()

	p.done = true
	p.b = v
	p.ex = e
}

// SetError indicates outward.Get() failed
func (p *outreacher[K]) SetError(err error) {
	defer p.cond.Broadcast()

	p.done = true
	p.err = err
}

// Wait patiently waits until the outreacher has finished its attempt
// to get the data
func (p *outreacher[K]) Wait() {
	p.count++
	for !p.done {
		p.cond.Wait()
	}
	p.count--
	p.Done()
}

// Done makes the [SingleFlight] parent forget about the block on this key
func (p *outreacher[K]) Done() bool {
	if p.count < 1 {
		delete(p.parent.getters, p.key)
		return true
	}
	return false
}
