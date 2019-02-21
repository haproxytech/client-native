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

type binds map[string]models.Bind

type bindCache struct {
	items map[string]map[string]binds
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

type stickRules map[int64]models.StickRule

type stickCache struct {
	items map[string]map[string]stickRules
	mu    sync.RWMutex
}

type tcpRequestRules map[int64]models.TCPRequestRule

type tcpRequestCache struct {
	items map[string]map[string]tcpRequestRules
	mu    sync.RWMutex
}

type tcpResponseRules map[int64]models.TCPResponseRule

type tcpResponseCache struct {
	items map[string]map[string]tcpResponseRules
	mu    sync.RWMutex
}

type logTargets map[int64]models.LogTarget

type logTargetCache struct {
	items map[string]map[string]logTargets
	mu    sync.RWMutex
}

type versionCache struct {
	items map[string]int64
	mu    sync.RWMutex
}

type Cache struct {
	enabled               bool
	Frontends             frontendCache
	Backends              backendCache
	Servers               serverCache
	Binds                 bindCache
	BackendSwitchingRules backendSwitchingCache
	Filters               filterCache
	HttpRequestRules      httpReqCache
	HttpResponseRules     httpResCache
	ServerSwitchingRules  serverSwitchingCache
	StickRules            stickCache
	TcpRequestRules       tcpRequestCache
	TcpResponseRules      tcpResponseCache
	LogTargets            logTargetCache
	Version               versionCache
}

func (c *Cache) Init(v int64) {
	c.enabled = true
	c.Version.items = make(map[string]int64)
	c.Version.Set(v, "")

	c.Frontends.items = make(map[string]map[string]models.Frontend)
	c.Backends.items = make(map[string]map[string]models.Backend)
	c.Binds.items = make(map[string]map[string]binds)
	c.Servers.items = make(map[string]map[string]servers)
	c.BackendSwitchingRules.items = make(map[string]map[string]backendSwitchingRules)
	c.Filters.items = make(map[string]map[string]filters)
	c.HttpRequestRules.items = make(map[string]map[string]httpReqRules)
	c.HttpResponseRules.items = make(map[string]map[string]httpResRules)
	c.ServerSwitchingRules.items = make(map[string]map[string]serverSwitchingRules)
	c.StickRules.items = make(map[string]map[string]stickRules)
	c.TcpRequestRules.items = make(map[string]map[string]tcpRequestRules)
	c.TcpResponseRules.items = make(map[string]map[string]tcpResponseRules)
	c.LogTargets.items = make(map[string]map[string]logTargets)
	c.InvalidateCache()
}

func (c *Cache) InitTransactionCache(t string, v int64) {
	c.Version.Set(v, t)
	c.Frontends.items[t] = make(map[string]models.Frontend)
	c.Backends.items[t] = make(map[string]models.Backend)
	c.Servers.items[t] = make(map[string]servers)
	c.Binds.items[t] = make(map[string]binds)
	c.BackendSwitchingRules.items[t] = make(map[string]backendSwitchingRules)
	c.Filters.items[t] = make(map[string]filters)
	c.HttpRequestRules.items[t] = make(map[string]httpReqRules)
	c.HttpResponseRules.items[t] = make(map[string]httpResRules)
	c.ServerSwitchingRules.items[t] = make(map[string]serverSwitchingRules)
	c.StickRules.items[t] = make(map[string]stickRules)
	c.TcpRequestRules.items[t] = make(map[string]tcpRequestRules)
	c.TcpResponseRules.items[t] = make(map[string]tcpResponseRules)
	c.LogTargets.items[t] = make(map[string]logTargets)
}

func (c *Cache) DeleteTransactionCache(t string) {
	delete(c.Version.items, t)
	delete(c.Frontends.items, t)
	delete(c.Backends.items, t)
	delete(c.Servers.items, t)
	delete(c.Binds.items, t)
	delete(c.BackendSwitchingRules.items, t)
	delete(c.Filters.items, t)
	delete(c.HttpRequestRules.items, t)
	delete(c.ServerSwitchingRules.items, t)
	delete(c.StickRules.items, t)
	delete(c.TcpRequestRules.items, t)
	delete(c.TcpResponseRules.items, t)
	delete(c.LogTargets.items, t)
}

func (c *Cache) DeleteFrontendCache(name, t string) {
	c.Frontends.delete(name, t)
	delete(c.Binds.items[t], name)
	delete(c.BackendSwitchingRules.items[t], name)
	delete(c.Filters.items[t], "frontend "+name)
	delete(c.HttpRequestRules.items[t], "frontend "+name)
	delete(c.HttpResponseRules.items[t], "frontend "+name)
	delete(c.TcpRequestRules.items[t], "frontend "+name)
	delete(c.LogTargets.items[t], "frontend "+name)
}

func (c *Cache) DeleteBackendCache(name, t string) {
	c.Backends.delete(name, t)
	delete(c.Servers.items[t], name)
	delete(c.Servers.items[t], "backend "+name)
	delete(c.HttpRequestRules.items[t], "backend "+name)
	delete(c.HttpResponseRules.items[t], "backend "+name)
	delete(c.ServerSwitchingRules.items[t], name)
	delete(c.StickRules.items[t], name)
	delete(c.TcpRequestRules.items[t], "backend "+name)
	delete(c.TcpResponseRules.items[t], "backend "+name)
	delete(c.LogTargets.items[t], "backend "+name)
}

func (c *Cache) InvalidateCache() {
	c.Frontends.Invalidate()
	c.Backends.Invalidate()
	c.Servers.Invalidate()
	c.Binds.Invalidate()
	c.BackendSwitchingRules.Invalidate()
	c.Filters.Invalidate()
	c.HttpRequestRules.Invalidate()
	c.HttpResponseRules.Invalidate()
	c.ServerSwitchingRules.Invalidate()
	c.StickRules.Invalidate()
	c.TcpRequestRules.Invalidate()
	c.TcpResponseRules.Invalidate()
	c.LogTargets.Invalidate()
}

func (c *Cache) Enabled() bool {
	return c.enabled
}

func (vc *versionCache) Get(transaction string) int64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	v, found := vc.items[transaction]
	if !found {
		return 0
	}
	return v
}

func (vc *versionCache) Set(v int64, transaction string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.items[transaction] = v
}

func (vc *versionCache) Increment() {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.items[""]++
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

func (lc *bindCache) Invalidate() {
	lc.mu.Lock()
	lc.items[""] = make(map[string]binds)
	lc.mu.Unlock()
}

func (lc *bindCache) Get(frontend, t string) (models.Binds, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	items, found := lc.items[t]
	if !found || len(items) == 0 {
		return nil, false
	}

	binds, found := items[frontend]
	if !found || len(binds) == 0 {
		return nil, false
	}

	ls := make([]*models.Bind, 0, len(binds))
	for _, l := range binds {
		lCopy := l
		ls = append(ls, &lCopy)
	}
	return ls, true
}

func (lc *bindCache) GetOne(name, frontend, t string) (*models.Bind, bool) {
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

func (lc *bindCache) Set(name, frontend, t string, l *models.Bind) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	_, found := lc.items[t]
	if !found {
		lc.items[t] = make(map[string]binds)
	}

	_, found = lc.items[t][frontend]
	if !found {
		lc.items[t][frontend] = make(binds)
	}
	lc.items[t][frontend][name] = *l
}

func (lc *bindCache) SetAll(frontend, t string, ls models.Binds) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	_, found := lc.items[t]
	if !found {
		lc.items[t] = make(map[string]binds)
	}

	_, found = lc.items[t][frontend]
	if !found {
		lc.items[t][frontend] = make(binds)
	}
	for _, l := range ls {
		lc.items[t][frontend][l.Name] = *l
	}
}

func (lc *bindCache) Delete(name, frontend, t string) {
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
		c.items[t][frontend][*r.ID] = *r
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
		c.items[t][parentType+" "+parent][*r.ID] = *r
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
		c.items[t][parentType+" "+parent][*r.ID] = *r
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
		c.items[t][parentType+" "+parent][*r.ID] = *r
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
		c.items[t][backend][*r.ID] = *r
	}
}

func (c *stickCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]stickRules)
	c.mu.Unlock()
}

func (c *stickCache) InvalidateBackend(t string, backend string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][backend] = make(stickRules)
	}
	c.mu.Unlock()
}

func (c *stickCache) Get(backend, t string) (models.StickRules, bool) {
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

	rs := make([]*models.StickRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *stickCache) GetOne(id int64, backend, t string) (*models.StickRule, bool) {
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

func (c *stickCache) Set(id int64, backend, t string, r *models.StickRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]stickRules)
	}
	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(stickRules)
	}
	c.items[t][backend][id] = *r
}

func (c *stickCache) SetAll(backend, t string, rs models.StickRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]stickRules)
	}

	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(stickRules)
	}
	for _, r := range rs {
		c.items[t][backend][*r.ID] = *r
	}
}

func (c *tcpRequestCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]tcpRequestRules)
	c.mu.Unlock()
}

func (c *tcpRequestCache) InvalidateParent(t, parent, parentType string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][parentType+" "+parent] = make(tcpRequestRules)
	}
	c.mu.Unlock()
}

func (c *tcpRequestCache) Get(parent, parentType, t string) (models.TCPRequestRules, bool) {
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

	rs := make([]*models.TCPRequestRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *tcpRequestCache) GetOne(id int64, parent, parentType, t string) (*models.TCPRequestRule, bool) {
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

func (c *tcpRequestCache) Set(id int64, parent, parentType, t string, r *models.TCPRequestRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpRequestRules)
	}
	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(tcpRequestRules)
	}
	c.items[t][parentType+" "+parent][id] = *r
}

func (c *tcpRequestCache) SetAll(parent, parentType, t string, rs models.TCPRequestRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpRequestRules)
	}

	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(tcpRequestRules)
	}
	for _, r := range rs {
		c.items[t][parentType+" "+parent][*r.ID] = *r
	}
}

func (c *tcpResponseCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]tcpResponseRules)
	c.mu.Unlock()
}

func (c *tcpResponseCache) InvalidateBackend(t string, backend string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][backend] = make(tcpResponseRules)
	}
	c.mu.Unlock()
}

func (c *tcpResponseCache) Get(backend, t string) (models.TCPResponseRules, bool) {
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

	rs := make([]*models.TCPResponseRule, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *tcpResponseCache) GetOne(id int64, backend, t string) (*models.TCPResponseRule, bool) {
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

func (c *tcpResponseCache) Set(id int64, backend, t string, r *models.TCPResponseRule) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpResponseRules)
	}
	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(tcpResponseRules)
	}
	c.items[t][backend][id] = *r
}

func (c *tcpResponseCache) SetAll(backend, t string, rs models.TCPResponseRules) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]tcpResponseRules)
	}

	_, found = c.items[t][backend]
	if !found {
		c.items[t][backend] = make(tcpResponseRules)
	}
	for _, r := range rs {
		c.items[t][backend][*r.ID] = *r
	}
}

func (c *logTargetCache) Invalidate() {
	c.mu.Lock()
	c.items[""] = make(map[string]logTargets)
	c.mu.Unlock()
}

func (c *logTargetCache) InvalidateParent(t, parent, parentType string) {
	c.mu.Lock()
	_, found := c.items[t]
	if found {
		c.items[t][parentType+" "+parent] = make(logTargets)
	}
	c.mu.Unlock()
}

func (c *logTargetCache) Get(parent, parentType, t string) (models.LogTargets, bool) {
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

	rs := make([]*models.LogTarget, 0, len(rules))
	for _, r := range rules {
		rCopy := r
		rs = append(rs, &rCopy)
	}
	return rs, true
}

func (c *logTargetCache) GetOne(id int64, parent, parentType, t string) (*models.LogTarget, bool) {
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

func (c *logTargetCache) Set(id int64, parent, parentType, t string, r *models.LogTarget) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]logTargets)
	}
	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(logTargets)
	}
	c.items[t][parentType+" "+parent][id] = *r
}

func (c *logTargetCache) SetAll(parent, parentType, t string, rs models.LogTargets) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[t]
	if !found {
		c.items[t] = make(map[string]logTargets)
	}

	_, found = c.items[t][parentType+" "+parent]
	if !found {
		c.items[t][parentType+" "+parent] = make(logTargets)
	}
	for _, r := range rs {
		c.items[t][parentType+" "+parent][*r.ID] = *r
	}
}
