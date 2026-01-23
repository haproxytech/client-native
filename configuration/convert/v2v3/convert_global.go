package v2v3

import (
	v2 "github.com/haproxytech/client-native/v5/models"
	v3 "github.com/haproxytech/client-native/v6/models"
)

func GlobalV2ToV3(v2g *v2.Global) (*v3.GlobalBase, error) { //nolint:maintidx
	var v3g *v3.GlobalBase
	// Global
	daemon := v2g.Daemon == "enabled"
	v3t, err := V2Tov3[v2.Global, v3.GlobalBase](v2g, "daemon")
	if err != nil {
		return nil, err
	}
	v3g = v3t

	// Fields restructured
	// PerformanceOptions
	performanceOptions, err := V2Tov3[v2.Global, v3.PerformanceOptions](v2g)
	if err != nil {
		return nil, err
	}
	if !performanceOptions.Equal(v3.PerformanceOptions{}) {
		v3g.PerformanceOptions = performanceOptions
	}

	// HTTPClientOptions
	httpClientOptions := &v3.HTTPClientOptions{
		ResolversDisabled: v2g.HttpclientResolversDisabled,
		ResolversID:       v2g.HttpclientResolversID,
		ResolversPrefer:   v2g.HttpclientResolversPrefer,
		Retries:           v2g.HttpclientRetries,
		SslCaFile:         v2g.HttpclientSslCaFile,
		SslVerify:         v2g.HttpclientSslVerify,
		TimeoutConnect:    v2g.HttpclientTimeoutConnect,
	}
	if !httpClientOptions.Equal(v3.HTTPClientOptions{}) {
		v3g.HTTPClientOptions = httpClientOptions
	}

	// TuneQuicOptions
	if v2g.TuneOptions != nil {
		tuneQuicOptions := &v3.TuneQuicOptions{
			FrontendConnTxBuffersLimit: v2g.TuneOptions.QuicFrontendConnTxBuffersLimit,
			FrontendMaxIdleTimeout:     v2g.TuneOptions.QuicFrontendMaxIdleTimeout,
			FrontendMaxStreamsBidi:     v2g.TuneOptions.QuicFrontendMaxStreamsBidi,
			MaxFrameLoss:               v2g.TuneOptions.QuicMaxFrameLoss,
			RetryThreshold:             v2g.TuneOptions.QuicRetryThreshold,
			SocketOwner:                v2g.TuneOptions.QuicSocketOwner,
			// ReorderRatio: not present in v2
			// ZeroCopyFwdSend: not present in v2
		}
		if !tuneQuicOptions.Equal(v3.TuneQuicOptions{}) {
			v3g.TuneQuicOptions = tuneQuicOptions
		}

		// TuneVarsOptions
		tuneVarsOptions := &v3.TuneVarsOptions{
			GlobalMaxSize: v2g.TuneOptions.VarsGlobalMaxSize,
			ProcMaxSize:   v2g.TuneOptions.VarsProcMaxSize,
			ReqresMaxSize: v2g.TuneOptions.VarsReqresMaxSize,
			SessMaxSize:   v2g.TuneOptions.VarsSessMaxSize,
			TxnMaxSize:    v2g.TuneOptions.VarsTxnMaxSize,
		}
		if !tuneVarsOptions.Equal(v3.TuneVarsOptions{}) {
			v3g.TuneVarsOptions = tuneVarsOptions
		}

		// TuneZlibOptions
		tuneZlibOptions := &v3.TuneZlibOptions{
			Memlevel:   v2g.TuneOptions.ZlibMemlevel,
			Windowsize: v2g.TuneOptions.ZlibWindowsize,
		}
		if !tuneZlibOptions.Equal(v3.TuneZlibOptions{}) {
			v3g.TuneZlibOptions = tuneZlibOptions
		}

		// TuneSslOptions
		tuneSslOptions := &v3.TuneSslOptions{
			Cachesize:          v2g.TuneOptions.SslCachesize,
			CtxCacheSize:       v2g.TuneOptions.SslCtxCacheSize,
			CaptureBufferSize:  v2g.TuneOptions.SslCaptureBufferSize,
			DefaultDhParam:     v2g.TuneOptions.SslDefaultDhParam,
			ForcePrivateCache:  v2g.TuneOptions.SslForcePrivateCache,
			Keylog:             v2g.TuneOptions.SslKeylog,
			Lifetime:           v2g.TuneOptions.SslLifetime,
			Maxrecord:          v2g.TuneOptions.SslMaxrecord,
			OcspUpdateMaxDelay: v2g.TuneOptions.SslOcspUpdateMaxDelay,
			OcspUpdateMinDelay: v2g.TuneOptions.SslOcspUpdateMinDelay,
		}
		if !tuneSslOptions.Equal(v3.TuneSslOptions{}) {
			v3g.TuneSslOptions = tuneSslOptions
		}

		// TuneLuaOptions
		tuneLuaOptions := &v3.TuneLuaOptions{
			BurstTimeout: v2g.TuneOptions.LuaBurstTimeout,
			ForcedYield:  v2g.TuneOptions.LuaForcedYield,
			LogLoggers:   v2g.TuneOptions.LuaLogLoggers,
			LogStderr:    v2g.TuneOptions.LuaLogStderr,
			// Maxmem:         was a boolean
			ServiceTimeout: v2g.TuneOptions.LuaServiceTimeout,
			SessionTimeout: v2g.TuneOptions.LuaSessionTimeout,
			TaskTimeout:    v2g.TuneOptions.LuaTaskTimeout,
		}
		if !tuneLuaOptions.Equal(v3.TuneLuaOptions{}) {
			v3g.TuneLuaOptions = tuneLuaOptions
		}

		// TuneBufOptions
		tuneBufferOptions := &v3.TuneBufferOptions{
			BuffersLimit:   v2g.TuneOptions.BuffersLimit,
			BuffersReserve: v2g.TuneOptions.BuffersReserve,
			Bufsize:        v2g.TuneOptions.Bufsize,
			Pipesize:       v2g.TuneOptions.Pipesize,
			RcvbufBackend:  v2g.TuneOptions.RcvbufBackend,
			RcvbufClient:   v2g.TuneOptions.RcvbufClient,
			RcvbufFrontend: v2g.TuneOptions.RcvbufFrontend,
			RcvbufServer:   v2g.TuneOptions.RcvbufServer,
			RecvEnough:     v2g.TuneOptions.RecvEnough,
			SndbufBackend:  v2g.TuneOptions.SndbufBackend,
			SndbufClient:   v2g.TuneOptions.SndbufClient,
			SndbufFrontend: v2g.TuneOptions.SndbufFrontend,
			SndbufServer:   v2g.TuneOptions.SndbufServer,
		}
		if !tuneBufferOptions.Equal(v3.TuneBufferOptions{}) {
			v3g.TuneBufferOptions = tuneBufferOptions
		}
	}

	// SslOptions
	v3engines, err := ListV2ToV3[v2.SslEngine, v3.SslEngine](v2g.SslEngines)
	if err != nil {
		return nil, err
	}
	sslOptions := &v3.SslOptions{
		SslEngines:                 v3engines,
		CaBase:                     v2g.CaBase,
		CrtBase:                    v2g.CrtBase,
		DefaultBindCiphers:         v2g.SslDefaultBindCiphers,
		DefaultBindCiphersuites:    v2g.SslDefaultBindCiphersuites,
		DefaultBindClientSigalgs:   v2g.SslDefaultBindClientSigalgs,
		DefaultBindCurves:          v2g.SslDefaultBindCurves,
		DefaultBindOptions:         v2g.SslDefaultBindOptions,
		DefaultBindSigalgs:         v2g.SslDefaultBindSigalgs,
		DefaultServerCiphers:       v2g.SslDefaultServerCiphers,
		DefaultServerCiphersuites:  v2g.SslDefaultServerCiphersuites,
		DefaultServerClientSigalgs: v2g.SslDefaultServerClientSigalgs,
		DefaultServerCurves:        v2g.SslDefaultServerCurves,
		DefaultServerOptions:       v2g.SslDefaultServerOptions,
		DefaultServerSigalgs:       v2g.SslDefaultServerSigalgs,
		DhParamFile:                v2g.SslDhParamFile,
		IssuersChainPath:           v2g.IssuersChainPath,
		LoadExtraFiles:             v2g.SslLoadExtraFiles,
		Maxsslconn:                 v2g.Maxsslconn,
		Maxsslrate:                 v2g.Maxsslrate,
		ModeAsync:                  v2g.SslModeAsync,
		Propquery:                  v2g.SslPropquery,
		Provider:                   v2g.SslProvider,
		ProviderPath:               v2g.SslProviderPath,
		// SecurityLevel: not present in v2
		ServerVerify:     v2g.SslServerVerify,
		SkipSelfIssuedCa: v2g.SslSkipSelfIssuedCa,
	}
	if !sslOptions.Equal(v3.SslOptions{}) {
		v3g.SslOptions = sslOptions
	}

	// EnvironmentOptions
	v3preset, err := ListV2ToV3[v2.PresetEnv, v3.PresetEnv](v2g.PresetEnvs)
	if err != nil {
		return nil, err
	}
	v3set, err := ListV2ToV3[v2.SetEnv, v3.SetEnv](v2g.SetEnvs)
	if err != nil {
		return nil, err
	}
	envOptions := &v3.EnvironmentOptions{
		PresetEnvs: v3preset,
		SetEnvs:    v3set,
		Resetenv:   v2g.Resetenv,
		Unsetenv:   v2g.Unsetenv,
	}
	if !envOptions.Equal(v3.EnvironmentOptions{}) {
		v3g.EnvironmentOptions = envOptions
	}

	// DebugOptions
	debugOptions := &v3.DebugOptions{
		Anonkey:     v2g.Anonkey,
		Quiet:       v2g.Quiet,
		ZeroWarning: v2g.ZeroWarning,
	}
	if !debugOptions.Equal(v3.DebugOptions{}) {
		v3g.DebugOptions = debugOptions
	}

	// LuaOptions
	luaLoad, err := ListV2ToV3[v2.LuaLoad, v3.LuaLoad](v2g.LuaLoads)
	if err != nil {
		return nil, err
	}
	luaPrependPath, err := ListV2ToV3[v2.LuaPrependPath, v3.LuaPrependPath](v2g.LuaPrependPath)
	if err != nil {
		return nil, err
	}
	luaOptions := &v3.LuaOptions{
		LoadPerThread: v2g.LuaLoadPerThread,
		Loads:         luaLoad,
		PrependPath:   luaPrependPath,
	}
	if !luaOptions.Equal(v3.LuaOptions{}) {
		v3g.LuaOptions = luaOptions
	}

	v3g.Daemon = daemon
	return v3g, nil
}
