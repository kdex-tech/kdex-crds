# API Reference

## Packages
- [kdex.dev/v1alpha1](#kdexdevv1alpha1)


## kdex.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group.

### Resource Types
- [KDexApp](#kdexapp)
- [KDexAppList](#kdexapplist)
- [KDexHost](#kdexhost)
- [KDexHostList](#kdexhostlist)
- [KDexPageArchetype](#kdexpagearchetype)
- [KDexPageArchetypeList](#kdexpagearchetypelist)
- [KDexPageBinding](#kdexpagebinding)
- [KDexPageBindingList](#kdexpagebindinglist)
- [KDexPageFooter](#kdexpagefooter)
- [KDexPageFooterList](#kdexpagefooterlist)
- [KDexPageHeader](#kdexpageheader)
- [KDexPageHeaderList](#kdexpageheaderlist)
- [KDexPageNavigation](#kdexpagenavigation)
- [KDexPageNavigationList](#kdexpagenavigationlist)
- [KDexRenderPage](#kdexrenderpage)
- [KDexRenderPageList](#kdexrenderpagelist)
- [KDexTranslation](#kdextranslation)
- [KDexTranslationList](#kdextranslationlist)
- [KDExTheme](#kdextheme)
- [KDExThemeList](#kdexthemelist)



#### AppPolicy

_Underlying type:_ _string_

AppPolicy defines the policy for apps.

_Validation:_
- Enum: [Strict NonStrict]

_Appears in:_
- [KDexHostSpec](#kdexhostspec)

| Field | Description |
| --- | --- |
| `Strict` | StrictAppPolicy means that apps may not embed JavaScript dependencies.<br /> |
| `NonStrict` | NonStrictAppPolicy means that apps may embed JavaScript dependencies.<br /> |






#### ContentEntry







_Appears in:_
- [KDexPageBindingSpec](#kdexpagebindingspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `appRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | appRef is a reference to the KDexApp to include in this binding. |  |  |
| `customElementName` _string_ | customElementName is the name of the KDexApp custom element to render in the specified slot (if present in the template). |  |  |
| `rawHTML` _string_ | rawHTML is a raw HTML string to be rendered in the specified slot (if present in the template). |  |  |
| `slot` _string_ | slot is the name of the App slot to which this entry will be bound. If omitted, the slot used will be `main`. No more than one entry can be bound to a slot. |  |  |


#### CustomElement



CustomElement defines a custom element exposed by a micro-frontend application.



_Appears in:_
- [KDexAppSpec](#kdexappspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `description` _string_ | description of the custom element. |  |  |
| `name` _string_ | name of the custom element. |  | Required: \{\} <br /> |


#### KDexApp



KDexApp is the Schema for the kdexapps API



_Appears in:_
- [KDexAppList](#kdexapplist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexApp` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexAppSpec](#kdexappspec)_ | spec defines the desired state of KDexApp |  | Required: \{\} <br /> |


#### KDexAppList



KDexAppList contains a list of KDexApp





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexAppList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexApp](#kdexapp) array_ |  |  |  |


#### KDexAppSpec



KDexAppSpec defines the desired state of KDexApp



_Appears in:_
- [KDexApp](#kdexapp)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `customElements` _[CustomElement](#customelement) array_ | customElements is a list of custom elements implemented by the micro-frontend application. |  |  |
| `packageReference` _[PackageReference](#packagereference)_ | packageReference specifies the name and version of an NPM package that contains the micro-frontend application. The package must have a package.json that contains ES modules. |  | Required: \{\} <br /> |




#### KDexHost



KDexHost is the Schema for the kdexhosts API



_Appears in:_
- [KDexHostList](#kdexhostlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexHost` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexHostSpec](#kdexhostspec)_ | spec defines the desired state of KDexHost |  |  |


#### KDexHostList



KDexHostList contains a list of KDexHost





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexHostList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexHost](#kdexhost) array_ |  |  |  |


#### KDexHostSpec



KDexHostSpec defines the desired state of KDexHost



_Appears in:_
- [KDexHost](#kdexhost)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `appPolicy` _[AppPolicy](#apppolicy)_ | AppPolicy defines the policy for apps.<br />When the strict policy is enabled, an app may not embed JavaScript dependencies.<br />Validation of the application source code will fail if dependencies are not fully externalized.<br />A Host which defines the `script` app policy must not accept apps which do not comply.<br />While a non-strict Host may accept both strict and non-strict apps. |  | Enum: [Strict NonStrict] <br />Required: \{\} <br /> |
| `baseMeta` _string_ | baseMeta is a string containing a base set of meta tags to use on every page rendered for the host. |  | MinLength: 5 <br /> |
| `defaultLang` _string_ | defaultLang is a string containing a BCP 47 language tag.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.<br />When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'. |  |  |
| `defaultThemeRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultThemeRef is a reference to the default theme that should apply to all pages bound to this host. |  |  |
| `domains` _string array_ | domains are the names by which this host is addressed. The first domain listed is the preferred domain. The domains may contain wildcard prefix in the form '*.'. Longest match always wins. |  | MinItems: 1 <br />Required: \{\} <br /> |
| `organization` _string_ | organization is the name of the Organization. |  | MinLength: 5 <br />Required: \{\} <br /> |




#### KDexPageArchetype



KDexPageArchetype is the Schema for the kdexpagearchetypes API



_Appears in:_
- [KDexPageArchetypeList](#kdexpagearchetypelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageArchetype` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexPageArchetypeSpec](#kdexpagearchetypespec)_ | spec defines the desired state of KDexPageArchetype |  | Required: \{\} <br /> |


#### KDexPageArchetypeList



KDexPageArchetypeList contains a list of KDexPageArchetype





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageArchetypeList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexPageArchetype](#kdexpagearchetype) array_ |  |  |  |


#### KDexPageArchetypeSpec



KDexPageArchetypeSpec defines the desired state of KDexPageArchetype



_Appears in:_
- [KDexPageArchetype](#kdexpagearchetype)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the structure of an HTML page. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `defaultFooterRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultFooterRef is an optional reference to a KDexPageFooter resource. If not specified, no footer will be displayed. Use the `.Footer` property to position its content in the template. |  |  |
| `defaultHeaderRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, no header will be displayed. Use the `.Header` property to position its content in the template. |  |  |
| `defaultMainNavigationRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultMainNavigationRef is an optional reference to a KDexPageNavigation resource. If not specified, no navigation will be displayed. Use the `.Navigation["main"]` property to position its content in the template. |  |  |
| `extraNavigations` _object (keys:string, values:[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core))_ | extraNavigations is an optional map of named navigation object references. Use `.Navigation["<name>"]` to position the named navigation's content in the template. |  |  |
| `overrideThemeRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideThemeRef is a reference to the theme that should apply to all pages that use this archetype. It overrides the default theme defined on the host. |  |  |




#### KDexPageBinding



KDexPageBinding is the Schema for the kdexpagebindings API



_Appears in:_
- [KDexPageBindingList](#kdexpagebindinglist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageBinding` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexPageBindingSpec](#kdexpagebindingspec)_ | spec defines the desired state of KDexPageBinding |  | Required: \{\} <br /> |


#### KDexPageBindingList



KDexPageBindingList contains a list of KDexPageBinding





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageBindingList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexPageBinding](#kdexpagebinding) array_ |  |  |  |


#### KDexPageBindingSpec



KDexPageBindingSpec defines the desired state of KDexPageBinding



_Appears in:_
- [KDexPageBinding](#kdexpagebinding)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `contentEntries` _[ContentEntry](#contententry) array_ | contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or KDexApp references. |  | MaxItems: 8 <br />MinItems: 1 <br />Required: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this binding is for. |  | Required: \{\} <br /> |
| `label` _string_ | label is the value used in menus and page titles before localization occurs (or when no translation exists for the current language). |  | Required: \{\} <br /> |
| `navigationHints` _[NavigationHints](#navigationhints)_ | navigationHints are optional navigation properties that if omitted result in the page being hidden from the navigation. |  |  |
| `overrideFooterRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideFooterRef is an optional reference to a KDexPageFooter resource. If not specified, the footer from the archetype will be used. |  |  |
| `overrideHeaderRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, the header from the archetype will be used. |  |  |
| `overrideMainNavigationRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideMainNavigationRef is an optional reference to a KDexPageNavigation resource. If not specified, the main navigation from the archetype will be used. |  |  |
| `pageArchetypeRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | pageArchetypeRef is a reference to the KDexPageArchetype that this binding is for. |  | Required: \{\} <br /> |
| `parentPageRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | parentPageRef is a reference to the KDexPageBinding bellow which this page will appear in the main navigation. If not set, the page will be placed in the top level of the navigation. |  |  |
| `basePath` _string_ | basePath is the shortest path by which the page may be accessed. It must not contain path parameters. This path will be used in site navigation. This path is subject to being prefixed for localization by `/\{l10n\}` and will be when the user selects a non-default language. |  | Pattern: `^/` <br />Required: \{\} <br /> |
| `patternPath` _string_ | patternPath, which must be prefixed by BasePath, is an extension of basePath that adds pattern matching as defined by https://pkg.go.dev/net/http#hdr-Patterns-ServeMux. This path is subject to being prefixed for localization by `/\{l10n\}` such as when the user selects a non-default language. |  |  |




#### KDexPageFooter



KDexPageFooter is the Schema for the kdexpagefooters API



_Appears in:_
- [KDexPageFooterList](#kdexpagefooterlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageFooter` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexPageFooterSpec](#kdexpagefooterspec)_ | spec defines the desired state of KDexPageFooter |  |  |


#### KDexPageFooterList



KDexPageFooterList contains a list of KDexPageFooter





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageFooterList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexPageFooter](#kdexpagefooter) array_ |  |  |  |


#### KDexPageFooterSpec



KDexPageFooterSpec defines the desired state of KDexPageFooter



_Appears in:_
- [KDexPageFooter](#kdexpagefooter)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the content of an App Server page footer section. Use the `.Footer` property to position its content in the template. |  | MinLength: 5 <br />Required: \{\} <br /> |




#### KDexPageHeader



KDexPageHeader is the Schema for the kdexpageheaders API



_Appears in:_
- [KDexPageHeaderList](#kdexpageheaderlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageHeader` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexPageHeaderSpec](#kdexpageheaderspec)_ | spec defines the desired state of KDexPageHeader |  | Required: \{\} <br /> |


#### KDexPageHeaderList



KDexPageHeaderList contains a list of KDexPageHeader





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageHeaderList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexPageHeader](#kdexpageheader) array_ |  |  |  |


#### KDexPageHeaderSpec



KDexPageHeaderSpec defines the desired state of KDexPageHeader



_Appears in:_
- [KDexPageHeader](#kdexpageheader)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the content of an App Server page header section. Use the `.Header` property to position its content in the template. |  | MinLength: 5 <br />Required: \{\} <br /> |




#### KDexPageNavigation



KDexPageNavigation is the Schema for the kdexpagenavigations API



_Appears in:_
- [KDexPageNavigationList](#kdexpagenavigationlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageNavigation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexPageNavigationSpec](#kdexpagenavigationspec)_ | spec defines the desired state of KDexPageNavigation |  | Required: \{\} <br /> |


#### KDexPageNavigationList



KDexPageNavigationList contains a list of KDexPageNavigation





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageNavigationList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexPageNavigation](#kdexpagenavigation) array_ |  |  |  |


#### KDexPageNavigationSpec



KDexPageNavigationSpec defines the desired state of KDexPageNavigation



_Appears in:_
- [KDexPageNavigation](#kdexpagenavigation)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the content of an App Server page navigation. Use the `.Navigation["<name>"]` property to position its content in the template. |  | MinLength: 5 <br />Required: \{\} <br /> |




#### KDexRenderPage



KDexRenderPage is the Schema for the kdexrenderpages API



_Appears in:_
- [KDexRenderPageList](#kdexrenderpagelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexRenderPage` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexRenderPageSpec](#kdexrenderpagespec)_ | spec defines the desired state of KDexRenderPage |  |  |


#### KDexRenderPageList



KDexRenderPageList contains a list of KDexRenderPage





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexRenderPageList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexRenderPage](#kdexrenderpage) array_ |  |  |  |


#### KDexRenderPageSpec



KDexRenderPageSpec defines the desired state of KDexRenderPage



_Appears in:_
- [KDexRenderPage](#kdexrenderpage)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this render page is for. |  | Required: \{\} <br /> |
| `navigationHints` _[NavigationHints](#navigationhints)_ | navigationHints are optional navigation properties that if omitted result in the page being hidden from the navigation. |  |  |
| `pageComponents` _[PageComponents](#pagecomponents)_ | pageComponents make up the elements of an HTML page that will be rendered by a web server. |  | Required: \{\} <br /> |
| `parentPageRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | parentPageRef is a reference to the KDexRenderPage bellow which this page will appear in the main navigation. If not set, the page will be placed in the top level of the navigation. |  |  |
| `basePath` _string_ | basePath is the shortest path by which the page may be accessed. It must not contain path parameters. This path will be used in site navigation. This path is subject to being prefixed for localization by `/\{l10n\}` and will be when the user selects a non-default language. |  | Pattern: `^/` <br />Required: \{\} <br /> |
| `patternPath` _string_ | patternPath, which must be prefixed by BasePath, is an extension of basePath that adds pattern matching as defined by https://pkg.go.dev/net/http#hdr-Patterns-ServeMux. This path is subject to being prefixed for localization by `/\{l10n\}` such as when the user selects a non-default language. |  |  |
| `themeRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | themeRef is a reference to the theme that will apply to this render page. |  |  |




#### KDexTranslation



KDexTranslation is the Schema for the kdextranslations API



_Appears in:_
- [KDexTranslationList](#kdextranslationlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexTranslation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexTranslationSpec](#kdextranslationspec)_ | spec defines the desired state of KDexTranslation |  |  |


#### KDexTranslationList



KDexTranslationList contains a list of KDexTranslation





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexTranslationList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexTranslation](#kdextranslation) array_ |  |  |  |


#### KDexTranslationSpec



KDexTranslationSpec defines the desired state of KDexTranslation



_Appears in:_
- [KDexTranslation](#kdextranslation)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this render page is for. |  | Required: \{\} <br /> |
| `translations` _[Translation](#translation) array_ | translations is an array of objects where each one specifies a language (lang) and a map (keysAndValues) consisting of key/value pairs. If the lang property is not unique in the array and its keysAndValues map contains the same keys, the last one takes precedence. |  | MinItems: 1 <br />Required: \{\} <br /> |




#### KDExTheme



KDExTheme is the Schema for the kdexthemes API



_Appears in:_
- [KDExThemeList](#kdexthemelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDExTheme` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDExThemeSpec](#kdexthemespec)_ | spec defines the desired state of KDExTheme |  |  |


#### KDExThemeList



KDExThemeList contains a list of KDExTheme





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDExThemeList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDExTheme](#kdextheme) array_ |  |  |  |


#### KDExThemeSpec



KDExThemeSpec defines the desired state of KDExTheme



_Appears in:_
- [KDExTheme](#kdextheme)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `styleItems` _[StyleItem](#styleitem) array_ | styleItems is a set of elements that define a portable set of design rules. They may contain URLs that point to resources hosted at some public address and/or they may contain the literal CSS. |  | MaxItems: 32 <br />MinItems: 1 <br />Required: \{\} <br /> |




#### NavigationHints







_Appears in:_
- [KDexPageBindingSpec](#kdexpagebindingspec)
- [KDexRenderPageSpec](#kdexrenderpagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `icon` _string_ | icon is the name of the icon to display next to the menu entry for this page. |  |  |
| `weight` _[Quantity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#quantity-resource-api)_ | weight is a property that influences the position of the page menu entry. Items at each level are sorted first by ascending weight and then ascending lexicographically. |  |  |


#### PackageReference



PackageReference specifies the name and version of an NPM package that contains the micro-frontend application.



_Appears in:_
- [KDexAppSpec](#kdexappspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `secretRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | secretRef is a reference to a secret containing authentication credentials for the NPM registry that holds the package. |  |  |
| `name` _string_ | name contains a scoped npm package name. |  | Required: \{\} <br /> |
| `version` _string_ | version contains a specific npm package version. |  | Required: \{\} <br /> |


#### PageComponents







_Appears in:_
- [KDexRenderPageSpec](#kdexrenderpagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `contents` _object (keys:string, values:string)_ |  |  |  |
| `footer` _string_ |  |  |  |
| `header` _string_ |  |  |  |
| `navigations` _object (keys:string, values:string)_ |  |  |  |
| `primaryTemplate` _string_ |  |  |  |
| `title` _string_ |  |  |  |


#### Paths







_Appears in:_
- [KDexPageBindingSpec](#kdexpagebindingspec)
- [KDexRenderPageSpec](#kdexrenderpagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `basePath` _string_ | basePath is the shortest path by which the page may be accessed. It must not contain path parameters. This path will be used in site navigation. This path is subject to being prefixed for localization by `/\{l10n\}` and will be when the user selects a non-default language. |  | Pattern: `^/` <br />Required: \{\} <br /> |
| `patternPath` _string_ | patternPath, which must be prefixed by BasePath, is an extension of basePath that adds pattern matching as defined by https://pkg.go.dev/net/http#hdr-Patterns-ServeMux. This path is subject to being prefixed for localization by `/\{l10n\}` such as when the user selects a non-default language. |  |  |


#### StyleItem







_Appears in:_
- [KDExThemeSpec](#kdexthemespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `attributes` _object (keys:string, values:string)_ | attributes are key/value pairs that will be added to the element [link\|style] when rendered. |  |  |
| `linkHref` _string_ | linkHref is the content of a <link> href attribute. |  |  |
| `style` _string_ | style is the text content to be added into a <script> element when rendered. |  |  |


#### Translation







_Appears in:_
- [KDexTranslationSpec](#kdextranslationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `keysAndValues` _object (keys:string, values:string)_ | keysAndValues is a map of key=/value pairs where the key is the identifier and the value is the translation of that key in the language specified by the lang property. |  | MinProperties: 1 <br />Required: \{\} <br /> |
| `lang` _string_ | lang is a string containing a BCP 47 language tag that identifies the set of translations.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag. |  | Required: \{\} <br /> |


