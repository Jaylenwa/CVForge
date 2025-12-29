package pdf

type ScreenshotRequest struct {
	URL              string         `json:"url"`
	HTML             string         `json:"html,omitempty"`
	Options          map[string]any `json:"options,omitempty"`
	EmulateMediaType string         `json:"emulateMediaType,omitempty"`
	GotoOptions      struct {
		Referer        string   `json:"referer,omitempty"`
		ReferrerPolicy string   `json:"referrerPolicy,omitempty"`
		Timeout        int      `json:"timeout,omitempty"`
		WaitUntil      []string `json:"waitUntil,omitempty"`
	} `json:"gotoOptions,omitempty"`
	WaitForSelector struct {
		Hidden   bool   `json:"hidden,omitempty"`
		Selector string `json:"selector,omitempty"`
		Timeout  int    `json:"timeout,omitempty"`
		Visible  bool   `json:"visible,omitempty"`
	} `json:"waitForSelector,omitempty"`
	WaitForTimeout      int               `json:"waitForTimeout,omitempty"`
	SetExtraHTTPHeaders map[string]string `json:"setExtraHTTPHeaders,omitempty"`
	BestAttempt         bool              `json:"bestAttempt,omitempty"`
	AddScriptTag        []struct {
		URL     string `json:"url,omitempty"`
		Path    string `json:"path,omitempty"`
		Content string `json:"content,omitempty"`
		Type    string `json:"type,omitempty"`
		ID      string `json:"id,omitempty"`
	} `json:"addScriptTag,omitempty"`
	AddStyleTag []struct {
		URL     string `json:"url,omitempty"`
		Path    string `json:"path,omitempty"`
		Content string `json:"content,omitempty"`
	} `json:"addStyleTag,omitempty"`
	Cookies []struct {
		Name         string `json:"name,omitempty"`
		Value        string `json:"value,omitempty"`
		URL          string `json:"url,omitempty"`
		Domain       string `json:"domain,omitempty"`
		Path         string `json:"path,omitempty"`
		Secure       bool   `json:"secure,omitempty"`
		HttpOnly     bool   `json:"httpOnly,omitempty"`
		SameSite     string `json:"sameSite,omitempty"`
		Expires      int64  `json:"expires,omitempty"`
		Priority     string `json:"priority,omitempty"`
		SameParty    bool   `json:"sameParty,omitempty"`
		SourceScheme string `json:"sourceScheme,omitempty"`
		PartitionKey struct {
			SourceOrigin         string `json:"sourceOrigin,omitempty"`
			HasCrossSiteAncestor bool   `json:"hasCrossSiteAncestor,omitempty"`
		} `json:"partitionKey,omitempty"`
	} `json:"cookies,omitempty"`
	RejectRequestPattern []string `json:"rejectRequestPattern,omitempty"`
	RejectResourceTypes  []string `json:"rejectResourceTypes,omitempty"`
	RequestInterceptors  []struct {
		Pattern  string `json:"pattern,omitempty"`
		Response struct {
			Headers     map[string]string `json:"headers,omitempty"`
			Status      int               `json:"status,omitempty"`
			ContentType string            `json:"contentType,omitempty"`
			Body        string            `json:"body,omitempty"`
		} `json:"response,omitempty"`
	} `json:"requestInterceptors,omitempty"`
	ScrollPage           bool   `json:"scrollPage,omitempty"`
	Selector             string `json:"selector,omitempty"`
	SetJavaScriptEnabled bool   `json:"setJavaScriptEnabled,omitempty"`
	UserAgent            *struct {
		UserAgent         string `json:"userAgent,omitempty"`
		Platform          string `json:"platform,omitempty"`
		UserAgentMetadata struct {
			Brands []struct {
				Brand   string `json:"brand,omitempty"`
				Version string `json:"version,omitempty"`
			} `json:"brands,omitempty"`
			FullVersionList []struct {
				Brand   string `json:"brand,omitempty"`
				Version string `json:"version,omitempty"`
			} `json:"fullVersionList,omitempty"`
			FullVersion     string   `json:"fullVersion,omitempty"`
			Platform        string   `json:"platform,omitempty"`
			PlatformVersion string   `json:"platformVersion,omitempty"`
			Architecture    string   `json:"architecture,omitempty"`
			Model           string   `json:"model,omitempty"`
			Mobile          bool     `json:"mobile,omitempty"`
			Bitness         string   `json:"bitness,omitempty"`
			Wow64           bool     `json:"wow64,omitempty"`
			FormFactors     []string `json:"formFactors,omitempty"`
		} `json:"userAgentMetadata,omitempty"`
	} `json:"userAgent,omitempty"`
	Viewport struct {
		Width             int     `json:"width,omitempty"`
		Height            int     `json:"height,omitempty"`
		DeviceScaleFactor float64 `json:"deviceScaleFactor,omitempty"`
		IsMobile          bool    `json:"isMobile,omitempty"`
		IsLandscape       bool    `json:"isLandscape,omitempty"`
		HasTouch          bool    `json:"hasTouch,omitempty"`
	} `json:"viewport,omitempty"`
}

type PDFRequest struct {
	URL              string         `json:"url"`
	HTML             string         `json:"html,omitempty"`
	Options          map[string]any `json:"options,omitempty"`
	EmulateMediaType string         `json:"emulateMediaType,omitempty"`
	GotoOptions      struct {
		Referer        string   `json:"referer,omitempty"`
		ReferrerPolicy string   `json:"referrerPolicy,omitempty"`
		Timeout        int      `json:"timeout,omitempty"`
		WaitUntil      []string `json:"waitUntil,omitempty"`
	} `json:"gotoOptions,omitempty"`
	WaitForSelector struct {
		Hidden   bool   `json:"hidden,omitempty"`
		Selector string `json:"selector,omitempty"`
		Timeout  int    `json:"timeout,omitempty"`
		Visible  bool   `json:"visible,omitempty"`
	} `json:"waitForSelector,omitempty"`
	SetExtraHTTPHeaders map[string]string `json:"setExtraHTTPHeaders,omitempty"`
	BestAttempt         bool              `json:"bestAttempt,omitempty"`
}
