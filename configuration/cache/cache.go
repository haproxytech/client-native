package cache

import (
	"sync"

	"github.com/haproxytech/models"
)

type frontendCache struct {
	items map[string]map[string]models.Frontend
	mu    sync.RWMutex
}

type backendCache struct {
	items map[string]map[string]models.Backend
	mu    sync.RWMutex
}

type servers map[string]models.Server

type serverCache struct {
	items map[string]map[string]servers
	mu    sync.RWMutex
}

type listeners map[string]models.Listener

type listenerCache struct {
	items map[string]map[string]listeners
	mu    sync.RWMutex
}

type backendSwitchingRules map[int64]models.BackendSwitchingRule

type backendSwitchingCache struct {
	items map[string]map[string]backendSwitchingRules
	mu    sync.RWMutex
}

type filters map[int64]models.Filter

type filterCache struct {
	items map[string]map[string]filters
	mu    sync.RWMutex
}

type httpReqRules map[int64]models.HTTPRequestRule

type httpReqCache struct {
	items map[string]map[string]httpReqRules
	mu    sync.RWMutex
}

type httpResRules map[int64]models.HTTPResponseRule

type httpResCache struct {
	items map[string]map[string]httpResRules
	mu    sync.RWMutex
}

type serverSwitchingRules map[int64]models.ServerSwitchingRule

type serverSwitchingCache struct {
	items map[string]map[string]serverSwitchingRules
	mu    sync.RWMutex
}

type stickRequestRules map[int64]models.StickRequestRule

type stickRequestCache struct {
	items map[string]map[string]stickRequestRules
	mu    sync.RWMutex
}

type stickResponseRules map[int64]models.StickResponseRule

type stickResponseCache struct {
	items map[string]map[string]stickResponseRules
	mu    sync.RWMutex
}

type tcpRules map[int64]models.TCPRule

type tcpConnCache struct {
	items map[string]map[string]tcpRules
	mu    sync.RWMutex
}

type tcpContReqCache struct {
	items map[string]map[string]tcpRules
	mu    sync.RWMutex
}

type tcpContResCache struct {
	items map[string]map[string]tcpRules
	mu    sync.RWMutex
}

type versionCache struct {
	v  int64
	mu sync.RWMutex
}

type Cache struct {
	enabled                 bool
	Frontends               frontendCache
	Backends                backendCache
	Servers                 serverCache
	Listeners               listenerCache
	BackendSwitchingRules   backendSwitchingCache
	Filters                 filterCache
	HttpRequestRules        httpReqCache
	HttpResponseRules       httpResCache
	ServerSwitchingRules    serverSwitchingCache
	StickRequestRules       stickRequestCache
	StickResponseRules      stickResponseCache
	TcpConnectionRules      tcpConnCache
	TcpContentRequestRules  tcpContReqCache
	TcpContentResponseRules tcpContResCache
	Version                 versionCache
}

func (c *Cache) Init(v int64) {
	c.enabled = true
	c.Version.Set(v)

	c.Frontends.items = make(map[string]map[string]models.Frontend)
	c.Backends.items = make(map[string]map[string]models.Backend)
	c.Listeners.items = make(map[string]map[string]listeners)
	c.Servers.items = make(map[string]map[string]servers)
	c.BackendSwitchingRules.items = make(map[string]map[string]backendSwitchingRules)
	c.Filters.items = make(map[string]map[string]filters)
	c.HttpRequestRules.items = make(map[string]map[string]httpReqRules)
	c.HttpResponseRules.items = make(map[string]map[string]httpResRules)
	c.ServerSwitchingRules.items = make(map[string]map[string]serverSwitchingRules)
	c.StickRequestRules.items = make(map[string]map[string]stickRequestRules)
	c.StickResponseRules.items = make(map[string]map[string]stickResponseRules)
	c.TcpConnectionRules.items = make(map[string]map[string]tcpRules)
	c.TcpContentRequestRules.items = make(map[string]map[string]tcpRules)
	c.TcpContentResponseRules.items = make(map[string]map[string]tcpRules)
	c.InvalidateCache()
}

func (c *Cache) InitTransactionCache(t string) {
	c.Frontends.items[t] = make(map[string]models.Frontend)
	c.Backends.items[t] = make(map[string]models.Backend)
	c.Servers.items[t] = make(map[string]servers)
	c.Listeners.items[t] = make(map[string]listeners)
	c.BackendSwitchingRules.items[t] = make(map[string]backendSwitchingRules)
	c.Filters.items[t] = make(map[string]filters)
	c.HttpRequestRules.items[t] = make(map[string]httpReqRules)
	c.HttpResponseRules.items[t] = make(map[string]httpResRules)
	c.ServerSwitchingRules.items[t] = make(map[string]serverSwitchingRules)
	c.StickRequestRules.items[t] = make(map[string]stickRequestRules)
	c.StickResponseRules.items[t] = make(map[string]stickResponseRules)
	c.TcpConnectionRules.items[t] = make(map[string]tcpRules)
	c.TcpContentRequestRules.items[t] = make(map[string]tcpRules)
	c.TcpContentResponseRules.items[t] = make(map[string]tcpRules)
}

func (c *Cache) DeleteTransactionCache(t string) {
	delete(c.Frontends.items, t)
	delete(c.Backends.items, t)
	delete(c.Servers.items, t)
	delete(c.Listeners.items, t)
	delete(c.BackendSwitchingRules.items, t)
	delete(c.Filters.items, t)
	delete(c.HttpRequestRules.items, t)
	delete(c.ServerSwitchingRules.items, t)
	delete(c.StickRequestRules.items, t)
	delete(c.StickResponseRules.items, t)
	delete(c.TcpConnectionRules.items, t)
	delete(c.TcpContentRequestRules.items, t)
	delete(c.TcpContentResponseRules.items, t)
}

func (c *Cache) DeleteFrontendCache(name, t string) {
	c.Frontends.delete(name, t)
	delete(c.Listeners.items[t], name)
	delete(c.BackendSwitchingRules.items[t], name)
	delete(c.Filters.items[t], "frontend "+name)
	delete(c.HttpRequestRules.items[t], "frontend "+name)
	delete(c.HttpResponseRules.items[t], "frontend "+name)
	delete(c.TcpConnectionRules.items[t], name)
	delete(c.TcpContentRequestRules.items[t], "frontend "+name)
}

func (c *Cache) DeleteBackendCache(name, t string) {
	c.Backends.delete(name, t)
	delete(c.Servers.items[t], name)
	delete(c.Servers.items[t], "backend "+name)
	delete(c.HttpRequestRules.items[t], "backend "+name)
	delete(c.HttpResponseRules.items[t], "backend "+name)
	delete(c.ServerSwitchingRules.items[t], name)
	delete(c.StickRequestRules.items[t], name)
	delete(c.StickResponseRules.items[t], name)
	delete(c.TcpContentRequestRules.items[t], "backend "+name)
	delete(c.TcpContentResponseRules.items[t], "backend "+name)
}

func (c *Cache) InvalidateCache() {
	c.Frontends.Invalidate()
	c.Backends.Invalidate()
	c.Servers.Invalidate()
	c.Listeners.Invalidate()
	c.BackendSwitchingRules.Invalidate()
	c.Filters.Invalidate()
	c.HttpRequestRules.Invalidate()
	c.HttpResponseRules.Invalidate()
	c.ServerSwitchingRules.Invalidate()
	c.StickRequestRules.Invalidate()
	c.StickResponseRules.Invalidate()
	c.TcpConnectionRules.Invalidate()
	c.TcpContentRequestRules.Invalidate()
	c.TcpContentResponseRules.Invalidate()
}

func (c *Cache) Enabled() bool {
	return c.enabled
}

func (vc *versionCache) Get() int64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return vc.v
}

func (vc *versionCache) Set(v int64) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.v = v
}

func (vc *versionCache) Increment() {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.v++
}

func (fc *frontendCache) Invalidate() {
	fc.mu.Lock()
	fc.items[""] = make(map[string]models.Frontend)
	fc.mu.Unlock()
}

func (fc *frontendCache) Get(t string) ([]*models.Frontend, bool) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	items, found := fc.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	fs := make([]*models.Frontend, 0, len(items))
	for _, val := range items {
		fCopy := val
		fs = append(fs, &fCopy)
	}
	return fs, true
}

func (fc *frontendCache) GetOne(name, t string) (*models.Frontend, bool) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	items, found := fc.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	f, found := items[name]
	if !found {
		return nil, false
	}
	return &f, true
}

func (fc *frontendCache) Set(name, t string, f *models.Frontend) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	_, found := fc.items[t]
	if !found {
		fc.items[t] = make(map[string]models.Frontend)
	}

	fc.items[t][name] = *f
}

func (fc *frontendCache) SetAll(t string, fs models.Frontends) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.items[t] = make(map[string]models.Frontend)

	for _, f := range fs {
		fc.items[t][f.Name] = *f
	}
}

func (fc *frontendCache) delete(name, t string) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	delete(fc.items[t], name)
}

func (bc *backendCache) Invalidate() {
	bc.mu.Lock()
	bc.items[""] = make(map[string]models.Backend)
	bc.mu.Unlock()
}

func (bc *backendCache) Get(t string) ([]*models.Backend, bool) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	items, found := bc.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	bs := make([]*models.Backend, 0, len(items))
	for _, val := range items {
		bCopy := val
		bs = append(bs, &bCopy)
	}
	return bs, true
}

func (bc *backendCache) GetOne(name, t string) (*models.Backend, bool) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	items, found := bc.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	b, found := items[name]
	if !found {
		return nil, false
	}
	return &b, true
}

func (bc *backendCache) Set(name, t string, b *models.Backend) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	_, found := bc.items[t]
	if !found {
		bc.items[t] = make(map[string]models.Backend)
	}

	bc.items[t][name] = *b
}

func (bc *backendCache) SetAll(t string, bs models.Backends) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.items[t] = make(map[string]models.Backend)
	for _, b := range bs {
		bc.items[t][b.Name] = *b
	}
}

func (bc *backendCache) delete(name, t string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	delete(bc.items[t], name)
}

func (lc *listenerCache) Invalidate() {
	lc.mu.Lock()
	lc.items[""] = make(map[string]listeners)
	lc.mu.Unlock()
}

func (lc *listenerCache) Get(frontend, t string) (models.Listeners, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	items, found := lc.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	listeners, found := items[frontend]
	if !found || len(listeners) == 0 {
		return nil, false
	}

	ls := make([]*models.Listener, 0, len(listeners))
	for _, l := range listeners {
		lCopy := l
		ls = append(ls, &lCopy)
	}
	return ls, true
}

func (lc *listenerCache) GetOne(name, frontend, t string) (*models.Listener, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	items, found := lc.items[t]
	if !found {
		return nil, false
	}

	ls, found := items[frontend]
	if !found {
		return nil, false
	}

	l, found := ls[name]
	if !found {
		return nil, false
	}
	return &l, true
}

func (lc *listenerCache) Set(name, frontend, t string, l *models.Listener) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	_, found := lc.items[t]
	if !found {
		lc.items[t] = make(map[string]listeners)
	}

	_, found = lc.items[t][frontend]
	if !found {
		lc.items[t][frontend] = make(listeners)
	}
	lc.items[t][frontend][name] = *l
}

func (lc *listenerCache) SetAll(frontend, t string, ls models.Listeners) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	_, found := lc.items[t]
	if !found {
		lc.items[t] = make(map[string]listeners)
	}

	_, found = lc.items[t][frontend]
	if !found {
		lc.items[t][frontend] = make(listeners)
	}
	for _, l := range ls {
		lc.items[t][frontend][l.Name] = *l
	}
}

func (lc *listenerCache) Delete(name, frontend, t string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	delete(lc.items[t][frontend], name)
}

func (sc *serverCache) Invalidate() {
	sc.mu.Lock()
	sc.items[""] = make(map[string]servers)
	sc.mu.Unlock()
}

func (sc *serverCache) Get(backend, t string) (models.Servers, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	items, found := sc.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	servers, found := items[backend]
	if !found || len(servers) == 0 {
		return nil, false
	}

	ss := make([]*models.Server, 0, len(servers))
	for _, s := range servers {
		sCopy := s
		ss = append(ss, &sCopy)
	}
	return ss, true
}

func (sc *serverCache) GetOne(name, backend, t string) (*models.Server, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	items, found := sc.items[t]
	if !found {
		return nil, false
	}

	ss, found := items[backend]
	if !found {
		return nil, false
	}
	s, found := ss[name]
	if !found {
		return nil, false
	}
	return &s, true
}

func (sc *serverCache) Set(name, backend, t string, s *models.Server) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	_, found := sc.items[t]
	if !found {
		sc.items[t] = make(map[string]servers)
	}
	_, found = sc.items[t][backend]
	if !found {
		sc.items[t][backend] = make(servers)
	}
	sc.items[t][backend][name] = *s
}

func (sc *serverCache) SetAll(backend, t string, ss models.Servers) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	_, found := sc.items[t]
	if !found {
		sc.items[t] = make(map[string]servers)
	}

	_, found = sc.items[t][backend]
	if !found {
		sc.items[t][backend] = make(servers)
	}
	for _, s := range ss {
		sc.items[t][backend][s.Name] = *s
	}
}

func (sc *serverCache) Delete(name, backend, t string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.items[t][backend], name)
}

func (c *backendSwitchingCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]backendSwitchingRules)
	c.mu.Unlock()
}

func (c *backendSwitchingCache) InvalidateFrontend(t string, frontend string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][frontend] = make(backendSwitchingRules)
	}
	c.mu.Unlock()
}

func (c *backendSwitchingCache) Get(frontend, t string) (models.BackendSwitchingRules, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[frontend]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.BackendSwitchingRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *backendSwitchingCache) GetOne(id int64, frontend, t string) (*models.BackendSwitchingRule, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[frontend]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *backendSwitchingCache) Set(id int64, frontend, t string, r *models.BackendSwitchingRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]backendSwitchingRules)
	}
	_, found = c.items[t][frontend]
	if !found {
		c.items[t][frontend] = make(backendSwitchingRules)
	}
	c.items[t][frontend][id] = *r
}

func (c *backendSwitchingCache) SetAll(frontend, t string, rs models.BackendSwitchingRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]backendSwitchingRules)
	}

	_, found = c.items[t][frontend]
	if !found {
		c.items[t][frontend] = make(backendSwitchingRules)
	}
	for _, r := range rs {
		c.items[t][frontend][r.ID] = *r
	}
}

func (c *filterCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]filters)
	c.mu.Unlock()
}

func (c *filterCache) InvalidateParent(t, parent, parentType string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][parentType+" "+parent] = make(filters)
	}
	c.mu.Unlock()
}

func (c *filterCache) Get(parent, parentType, t string) (models.Filters, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[parentType+" "+parent]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.Filter, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *filterCache) GetOne(id int64, parent, parentType, t string) (*models.Filter, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[parentType+" "+parent]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *filterCache) Set(id int64, parent, parentType, t string, r *models.Filter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]filters)
	}
	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(filters)
	}
	c.items[t][parentType+" "+parent][id] = *r
}

func (c *filterCache) SetAll(parent, parentType, t string, rs models.Filters) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]filters)
	}

	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(filters)
	}
	for _, r := range rs {
		c.items[t][parentType+" "+parent][r.ID] = *r
	}
}

func (c *httpReqCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]httpReqRules)
	c.mu.Unlock()
}

func (c *httpReqCache) InvalidateParent(t, parent, parentType string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][parentType+" "+parent] = make(httpReqRules)
	}
	c.mu.Unlock()
}

func (c *httpReqCache) Get(parent, parentType, t string) (models.HTTPRequestRules, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[parentType+" "+parent]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.HTTPRequestRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *httpReqCache) GetOne(id int64, parent, parentType, t string) (*models.HTTPRequestRule, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[parentType+" "+parent]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *httpReqCache) Set(id int64, parent, parentType, t string, r *models.HTTPRequestRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]httpReqRules)
	}
	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(httpReqRules)
	}
	c.items[t][parentType+" "+parent][id] = *r
}

func (c *httpReqCache) SetAll(parent, parentType, t string, rs models.HTTPRequestRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]httpReqRules)
	}

	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(httpReqRules)
	}
	for _, r := range rs {
		c.items[t][parentType+" "+parent][r.ID] = *r
	}
}

func (c *httpResCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]httpResRules)
	c.mu.Unlock()
}

func (c *httpResCache) InvalidateParent(t, parent, parentType string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][parentType+" "+parent] = make(httpResRules)
	}
	c.mu.Unlock()
}

func (c *httpResCache) Get(parent, parentType, t string) (models.HTTPResponseRules, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[parentType+" "+parent]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.HTTPResponseRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *httpResCache) GetOne(id int64, parent, parentType, t string) (*models.HTTPResponseRule, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[parentType+" "+parent]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *httpResCache) Set(id int64, parent, parentType, t string, r *models.HTTPResponseRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]httpResRules)
	}
	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(httpResRules)
	}
	c.items[t][parentType+" "+parent][id] = *r
}

func (c *httpResCache) SetAll(parent, parentType, t string, rs models.HTTPResponseRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]httpResRules)
	}

	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(httpResRules)
	}
	for _, r := range rs {
		c.items[t][parentType+" "+parent][r.ID] = *r
	}
}

func (c *serverSwitchingCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]serverSwitchingRules)
	c.mu.Unlock()
}

func (c *serverSwitchingCache) InvalidateBackend(t string, backend string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][backend] = make(serverSwitchingRules)
	}
	c.mu.Unlock()
}

func (c *serverSwitchingCache) Get(backend, t string) (models.ServerSwitchingRules, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[backend]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.ServerSwitchingRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *serverSwitchingCache) GetOne(id int64, backend, t string) (*models.ServerSwitchingRule, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[backend]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *serverSwitchingCache) Set(id int64, backend, t string, r *models.ServerSwitchingRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]serverSwitchingRules)
	}
	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(serverSwitchingRules)
	}
	c.items[t][backend][id] = *r
}

func (c *serverSwitchingCache) SetAll(backend, t string, rs models.ServerSwitchingRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]serverSwitchingRules)
	}

	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(serverSwitchingRules)
	}
	for _, r := range rs {
		c.items[t][backend][r.ID] = *r
	}
}

func (c *stickRequestCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]stickRequestRules)
	c.mu.Unlock()
}

func (c *stickRequestCache) InvalidateBackend(t string, backend string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][backend] = make(stickRequestRules)
	}
	c.mu.Unlock()
}

func (c *stickRequestCache) Get(backend, t string) (models.StickRequestRules, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[backend]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.StickRequestRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *stickRequestCache) GetOne(id int64, backend, t string) (*models.StickRequestRule, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[backend]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *stickRequestCache) Set(id int64, backend, t string, r *models.StickRequestRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]stickRequestRules)
	}
	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(stickRequestRules)
	}
	c.items[t][backend][id] = *r
}

func (c *stickRequestCache) SetAll(backend, t string, rs models.StickRequestRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]stickRequestRules)
	}

	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(stickRequestRules)
	}
	for _, r := range rs {
		c.items[t][backend][r.ID] = *r
	}
}

func (c *stickResponseCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]stickResponseRules)
	c.mu.Unlock()
}

func (c *stickResponseCache) InvalidateBackend(t string, backend string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][backend] = make(stickResponseRules)
	}
	c.mu.Unlock()
}

func (c *stickResponseCache) Get(backend, t string) (models.StickResponseRules, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[backend]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.StickResponseRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *stickResponseCache) GetOne(id int64, backend, t string) (*models.StickResponseRule, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[backend]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *stickResponseCache) Set(id int64, backend, t string, r *models.StickResponseRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]stickResponseRules)
	}
	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(stickResponseRules)
	}
	c.items[t][backend][id] = *r
}

func (c *stickResponseCache) SetAll(backend, t string, rs models.StickResponseRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]stickResponseRules)
	}

	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(stickResponseRules)
	}
	for _, r := range rs {
		c.items[t][backend][r.ID] = *r
	}
}

func (c *tcpConnCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]tcpRules)
	c.mu.Unlock()
}

func (c *tcpConnCache) InvalidateFrontend(t string, frontend string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][frontend] = make(tcpRules)
	}
	c.mu.Unlock()
}

func (c *tcpConnCache) Get(frontend, t string) (models.TCPRules, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[frontend]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.TCPRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *tcpConnCache) GetOne(id int64, frontend, t string) (*models.TCPRule, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[frontend]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *tcpConnCache) Set(id int64, frontend, t string, r *models.TCPRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpRules)
	}
	_, found = c.items[t][frontend]
	if !found {
		c.items[t][frontend] = make(tcpRules)
	}
	c.items[t][frontend][id] = *r
}

func (c *tcpConnCache) SetAll(frontend, t string, rs models.TCPRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpRules)
	}

	_, found = c.items[t][frontend]
	if !found {
		c.items[t][frontend] = make(tcpRules)
	}
	for _, r := range rs {
		c.items[t][frontend][r.ID] = *r
	}
}

func (c *tcpContReqCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]tcpRules)
	c.mu.Unlock()
}

func (c *tcpContReqCache) InvalidateParent(t, parent, parentType string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][parentType+" "+parent] = make(tcpRules)
	}
	c.mu.Unlock()
}

func (c *tcpContReqCache) Get(parent, parentType, t string) (models.TCPRules, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[parentType+" "+parent]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.TCPRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *tcpContReqCache) GetOne(id int64, parent, parentType, t string) (*models.TCPRule, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[parentType+" "+parent]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *tcpContReqCache) Set(id int64, parent, parentType, t string, r *models.TCPRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpRules)
	}
	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(tcpRules)
	}
	c.items[t][parentType+" "+parent][id] = *r
}

func (c *tcpContReqCache) SetAll(parent, parentType, t string, rs models.TCPRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpRules)
	}

	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(tcpRules)
	}
	for _, r := range rs {
		c.items[t][parentType+" "+parent][r.ID] = *r
	}
}

func (c *tcpContResCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]tcpRules)
	c.mu.Unlock()
}

func (c *tcpContResCache) InvalidateBackend(t string, backend string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][backend] = make(tcpRules)
	}
	c.mu.Unlock()
}

func (c *tcpContResCache) Get(backend, t string) (models.TCPRules, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	rules, found := items[backend]
	if !found || len(rules) == 0 {
		return nil, false
	}

	rs := make([]*models.TCPRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *tcpContResCache) GetOne(id int64, backend, t string) (*models.TCPRule, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items, found := c.items[t]
	if !found {
		return nil, false
	}

	rs, found := items[backend]
	if !found {
		return nil, false
	}

	r, found := rs[id]
	if !found {
		return nil, false
	}
	return &r, true
}

func (c *tcpContResCache) Set(id int64, backend, t string, r *models.TCPRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpRules)
	}
	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(tcpRules)
	}
	c.items[t][backend][id] = *r
}

func (c *tcpContResCache) SetAll(backend, t string, rs models.TCPRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpRules)
	}

	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(tcpRules)
	}
	for _, r := range rs {
		c.items[t][backend][r.ID] = *r
	}
}
