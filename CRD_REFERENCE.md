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
- [KDexScriptLibrary](#kdexscriptlibrary)
- [KDexScriptLibraryList](#kdexscriptlibrarylist)
- [KDexTheme](#kdextheme)
- [KDexThemeList](#kdexthemelist)
- [KDexTranslation](#kdextranslation)
- [KDexTranslationList](#kdextranslationlist)



#### Asset







_Appears in:_
- [KDexThemeSpec](#kdexthemespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `attributes` _object (keys:string, values:string)_ | attributes are key/value pairs that will be added to the element [link\|style] as attributes when rendered. |  |  |
| `linkHref` _string_ | linkHref is the content of a <link> href attribute. The URL may be absolute with protocol and host or it must be prefixed by the RoutePath of the theme. |  |  |
| `style` _string_ | style is the text content to be added into a <style> element when rendered. |  |  |




#### ConditionFields







_Appears in:_
- [ConditionArgs](#conditionargs)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `Status` _[ConditionStatus](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#conditionstatus-v1-meta)_ |  |  |  |
| `Reason` _[ConditionReason](#conditionreason)_ |  |  |  |
| `Message` _string_ |  |  |  |


#### ConditionReason

_Underlying type:_ _string_

ConditionReason is the reason for the condition's last transition.



_Appears in:_
- [ConditionFields](#conditionfields)

| Field | Description |
| --- | --- |
| `ReconcileError` | ConditionReasonReconcileError is the reason for a failed reconciliation.<br /> |
| `Reconciling` | ConditionReasonReconciling is the reason for a reconciling reconciliation.<br /> |
| `ReconcileSuccess` | ConditionReasonReconcileSuccess is the reason for a successful reconciliation.<br /> |




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



KDexApp is the Schema for the kdexapps API.

A KDexApp is the embodiment of an "Application" within the "KDex Cloud Native Application Server" model. KDexApp is
the resource developers implement to extend to user interface with a new feature. The implementations are Web
Component based and the packaging follows the NPM packaging model the contents of which are ES modules. There are no
container images to build. Merely package the application code and publish it to an NPM compatible repository,
configure the KDexApp with the necessary metadata and deploy to Kubernetes. The app can then be consumed and composed
by KDexPageBindings to produce actual user experiences.



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
| `customElements` _[CustomElement](#customelement) array_ | customElements is a list of custom elements implemented by the micro-frontend application. |  | MaxItems: 32 <br />MinItems: 1 <br /> |
| `packageReference` _[PackageReference](#packagereference)_ | packageReference specifies the name and version of an NPM package that contains the micro-frontend application. The package.json must describe an ES module. |  | Required: \{\} <br /> |
| `scripts` _[Script](#script) array_ | scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents. |  | MaxItems: 32 <br /> |




#### KDexHost



KDexHost is the Schema for the kdexhosts API

A KDexHost is the central actor in the "KDex Cloud Native Application Server" model. It specifies the basic metadata
that defines a web property; a set of domain names, TLS certificates, routing strategy and so on. From this central
point a distinct web property is establish to which are bound KDexPageBindings (i.e. web pages) that provide the web
properties content in the form of either raw HTML content or applications from KDexApps.s



_Appears in:_
- [KDexHostList](#kdexhostlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexHost` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexHostSpec](#kdexhostspec)_ | spec defines the desired state of KDexHost |  | Required: \{\} <br /> |


#### KDexHostList



KDexHostList contains a list of KDexHost





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexHostList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexHost](#kdexhost) array_ |  |  |  |


#### KDexHostSpec







_Appears in:_
- [KDexHost](#kdexhost)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `baseMeta` _string_ | baseMeta is a string containing a base set of meta tags to use on every page rendered for the host. |  | MinLength: 5 <br /> |
| `brandName` _string_ | brandName is the name used when rendering pages belonging to the host. For example, it may be used as alt text for the logo displayed in the page header. |  | Required: \{\} <br /> |
| `defaultLang` _string_ | defaultLang is a string containing a BCP 47 language tag.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.<br />When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'. |  |  |
| `defaultThemeRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultThemeRef is a reference to the theme that should apply to all pages bound to this host unless overridden. |  |  |
| `modulePolicy` _[ModulePolicy](#modulepolicy)_ | modulePolicy defines the policy for JavaScript references in KDexApp, KDexTheme and KDexScriptLibrary resources. When not specified the policy is Strict<br />A Host must not accept JavaScript references which do not comply with the specified policy. | Strict | Enum: [ExternalDependencies Loose ModulesRequired Strict] <br /> |
| `organization` _string_ | organization is the name of the Organization to which the host belongs. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `routing` _[Routing](#routing)_ | routing defines the desired routing configuration for the host. |  | Required: \{\} <br /> |
| `scriptLibraryRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  |  |




#### KDexPageArchetype



KDexPageArchetype is the Schema for the kdexpagearchetypes API

A KDexPageArchetype defines a reusable archetype from which web pages can be derived. When creating a KDexPageBinding
(i.e. a web page) a developer states which archetype is to be used. This allows the structure to be decoupled from
the content.



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
| `scriptLibraryRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  |  |




#### KDexPageBinding



KDexPageBinding is the Schema for the kdexpagebindings API

A KDexPageBinding defines a web page under a KDexHost. It brings together various reusable components like
KDexPageArchetype, KDexPageFooter, KDexPageHeader, KDexPageNavigation, KDexScriptLibrary, KDexTheme and content
components like raw HTML or KDexApps and KDexTranslations to produce internationalized, rendered HTML pages.



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
| `scriptLibraryRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  |  |




#### KDexPageFooter



KDexPageFooter is the Schema for the kdexpagefooters API

A KDexPageFooter is a reusable footer component for composing KDexPageBindings. It can specify a content template and
an associated KDexScriptLibrary for driving imperative logic that might be necessary to implement the footer.



_Appears in:_
- [KDexPageFooterList](#kdexpagefooterlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexPageFooter` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexPageFooterSpec](#kdexpagefooterspec)_ | spec defines the desired state of KDexPageFooter |  | Required: \{\} <br /> |


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
| `scriptLibraryRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  |  |




#### KDexPageHeader



KDexPageHeader is the Schema for the kdexpageheaders API

A KDexPageHeader is a reusable header component for composing KDexPageBindings. It can specify a content template and
an associated KDexScriptLibrary for driving imperative logic that might be necessary to implement the header.



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
| `scriptLibraryRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  |  |




#### KDexPageNavigation



KDexPageNavigation is the Schema for the kdexpagenavigations API

A KDexPageNavigation is a reusable navigation component for composing KDexPageBindings. It can specify a content
template and an associated KDexScriptLibrary for driving imperative logic that might be necessary to implement the
navigation.



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
| `scriptLibraryRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  |  |




#### KDexRenderPage



KDexRenderPage is the Schema for the kdexrenderpages API.
It is an internal resource created and managed by a controller that processes KDexPageBinding resources.
It is not intended for users to create or manage directly.



_Appears in:_
- [KDexRenderPageList](#kdexrenderpagelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexRenderPage` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexRenderPageSpec](#kdexrenderpagespec)_ | spec defines the desired state of KDexRenderPage |  | Required: \{\} <br /> |


#### KDexRenderPageList



KDexRenderPageList contains a list of KDexRenderPage





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexRenderPageList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexRenderPage](#kdexrenderpage) array_ |  |  |  |


#### KDexRenderPageSpec



KDexRenderPageSpec defines the desired state of KDexRenderPage.
KDexRenderPage is an internal resource created and managed by a controller that processes KDexPageBinding resources.
It is not intended for users to create or manage directly.



_Appears in:_
- [KDexRenderPage](#kdexrenderpage)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this render page is for. |  | Required: \{\} <br /> |
| `navigationHints` _[NavigationHints](#navigationhints)_ | navigationHints are optional navigation properties that if omitted result in the page being hidden from the navigation. |  |  |
| `pageComponents` _[PageComponents](#pagecomponents)_ | pageComponents make up the elements of an HTML page that will be rendered by a web server. |  | Required: \{\} <br /> |
| `parentPageRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | parentPageRef is a reference to the KDexRenderPage bellow which this page will appear in the main navigation. If not set, the page will be placed in the top level of the navigation. |  |  |
| `scriptLibraryRefs` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | scriptLibraryRefs is an optional array of KDexScriptLibrary references. |  |  |
| `themeRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | themeRef is a reference to the theme that will apply to this render page. |  |  |
| `basePath` _string_ | basePath is the shortest path by which the page may be accessed. It must not contain path parameters. This path will be used in site navigation. This path is subject to being prefixed for localization by `/\{l10n\}` and will be when the user selects a non-default language. |  | Pattern: `^/` <br />Required: \{\} <br /> |
| `patternPath` _string_ | patternPath, which must be prefixed by BasePath, is an extension of basePath that adds pattern matching as defined by https://pkg.go.dev/net/http#hdr-Patterns-ServeMux. This path is subject to being prefixed for localization by `/\{l10n\}` such as when the user selects a non-default language. |  |  |




#### KDexScriptLibrary



KDexScriptLibrary is the Schema for the kdexscriptlibraries API

A KDexScriptLibrary is a reusable collection of JavaScript for powering the imperative aspects of KDexPageBindings.
Most other components of the model are able to reference KDexScriptLibrary as well in order to encapsulate component
specific logic.



_Appears in:_
- [KDexScriptLibraryList](#kdexscriptlibrarylist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexScriptLibrary` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexScriptLibrarySpec](#kdexscriptlibraryspec)_ | spec defines the desired state of KDexScriptLibrary |  |  |


#### KDexScriptLibraryList



KDexScriptLibraryList contains a list of KDexScriptLibrary





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexScriptLibraryList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexScriptLibrary](#kdexscriptlibrary) array_ |  |  |  |


#### KDexScriptLibrarySpec



KDexScriptLibrarySpec defines the desired state of KDexScriptLibrary



_Appears in:_
- [KDexScriptLibrary](#kdexscriptlibrary)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `packageReference` _[PackageReference](#packagereference)_ | packageReference specifies the name and version of an NPM package that contains the script. The package.json must describe an ES module. |  |  |
| `scripts` _[Script](#script) array_ | scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents. |  | MaxItems: 32 <br /> |




#### KDexTheme



KDexTheme is the Schema for the kdexthemes API

A KDexTheme is a reusable collection of design styles and associated digital assets necessary for providing the
visual aspects of KDexPageBindings decoupling appearance from structure and content.



_Appears in:_
- [KDexThemeList](#kdexthemelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexTheme` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexThemeSpec](#kdexthemespec)_ | spec defines the desired state of KDexTheme |  | Required: \{\} <br /> |


#### KDexThemeList



KDexThemeList contains a list of KDexTheme





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexThemeList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexTheme](#kdextheme) array_ |  |  |  |


#### KDexThemeSpec



KDexThemeSpec defines the desired state of KDexTheme



_Appears in:_
- [KDexTheme](#kdextheme)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `assets` _[Asset](#asset) array_ | assets is a set of elements that define a portable set of design rules. |  | MaxItems: 32 <br />MinItems: 1 <br /> |
| `image` _string_ | image is the name of an OCI image that contains Theme resources.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  |  |
| `pullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  |  |
| `pullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | pullSecrets is an optional list of references to secrets in the same namespace to use for pulling the image. Also used for the webserver image if specified.<br />More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod |  |  |
| `routePath` _string_ | routePath is a prefix beginning with a forward slash (/) plus at least 1 additional character. KDexPageBindings associated with the KDexHost that have conflicting urls will be rejected and marked as conflicting. |  | Pattern: `^/.+` <br /> |
| `scriptLibraryRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  |  |
| `webserver` _[KDexThemeWebServer](#kdexthemewebserver)_ | webserver defines the configuration for the theme webserver. |  |  |




#### KDexThemeWebServer



KDexThemeWebServer defines the desired state of the KDexTheme web server



_Appears in:_
- [KDexThemeSpec](#kdexthemespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | image is the name of webserver image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Required: \{\} <br /> |
| `pullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the webserver image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  |  |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  |  |


#### KDexTranslation



KDexTranslation is the Schema for the kdextranslations API

KDexTranslations allow KDexPageBindings to be internationalized by making translations available in as many languages
as necessary.



_Appears in:_
- [KDexTranslationList](#kdextranslationlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexTranslation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexTranslationSpec](#kdextranslationspec)_ | spec defines the desired state of KDexTranslation |  | Required: \{\} <br /> |


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




#### ModulePolicy

_Underlying type:_ _string_

ModulePolicy defines the policy for the use of JavaScript Modules.

_Validation:_
- Enum: [ExternalDependencies Loose ModulesRequired Strict]

_Appears in:_
- [KDexHostSpec](#kdexhostspec)

| Field | Description |
| --- | --- |
| `Loose` | LooseModulePolicy means that a) JavaScript references are not required to be JavaScript modules and b) JavaScript references may contain embed dependencies.<br /> |
| `ExternalDependencies` | ExternalDependenciesModulePolicy means that a) JavaScript references are not required to be JavaScript modules and b) JavaScript references may not contain embed dependencies.<br /> |
| `ModulesRequired` | ModulesRequiredModulePolicy means that a) JavaScript references are required to be JavaScript modules and b) JavaScript references may contain embed dependencies.<br /> |
| `Strict` | StrictModulePolicy means that a) JavaScript references are required to be JavaScript modules and b) JavaScript references may not contain embed dependencies.<br /> |


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
- [KDexScriptLibrarySpec](#kdexscriptlibraryspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `exportMapping` _string_ | exportMapping is a mapping of the module's exports that will be used when the module import is written. e.g. `import [exportMapping] from [module_name];`. If exportMapping is not provided the module will be written as `import [module_name];` |  |  |
| `name` _string_ | name contains a scoped npm package name. |  | Required: \{\} <br /> |
| `secretRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | secretRef is a reference to a secret containing authentication credentials for the NPM registry that holds the package. |  |  |
| `version` _string_ | version contains a specific npm package version. |  | Required: \{\} <br /> |


#### PageComponents



PageComponents make up the elements of an HTML page that will be rendered by a web server.
It is an internal resource created and managed by a controller that processes KDexPageBinding resources.
It is not intended for users to create or manage directly.



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


#### Routing



Routing defines the desired routing configuration for the host.



_Appears in:_
- [KDexHostSpec](#kdexhostspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `domains` _string array_ | domains are the names by which this host is addressed. The first domain listed is the preferred domain. The domains may contain wildcard prefix in the form '*.'. Longest match always wins. |  | MinItems: 1 <br />Required: \{\} <br /> |
| `ingressClassName` _string_ | ingressClassName is the name of an IngressClass cluster resource. Ingress<br />controller implementations use this field to know whether they should be<br />serving this Ingress resource, by a transitive connection<br />(controller -> IngressClass -> Ingress resource). Although the<br />`kubernetes.io/ingress.class` annotation (simple constant name) was never<br />formally defined, it was widely supported by Ingress controllers to create<br />a direct binding between Ingress controller and Ingress resources. Newly<br />created Ingress resources should prefer using the field. However, even<br />though the annotation is officially deprecated, for backwards compatibility<br />reasons, ingress controllers should still honor that annotation if present. |  |  |
| `strategy` _[RoutingStrategy](#routingstrategy)_ | strategy is the routing strategy to use. If not specified Ingress is assumed. | Ingress | Enum: [Ingress HTTPRoute] <br /> |
| `tls` _[TLSSpec](#tlsspec)_ | tls is the TLS configuration for the host. |  |  |


#### RoutingStrategy

_Underlying type:_ _string_

RoutingStrategy defines the routing strategy to use.

_Validation:_
- Enum: [Ingress HTTPRoute]

_Appears in:_
- [Routing](#routing)

| Field | Description |
| --- | --- |
| `HTTPRoute` | HTTPRouteRoutingStrategy uses HTTPRoute to expose the host.<br /> |
| `Ingress` | IngressRoutingStrategy uses Ingress to expose the host.<br /> |


#### Script







_Appears in:_
- [KDexAppSpec](#kdexappspec)
- [KDexScriptLibrarySpec](#kdexscriptlibraryspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `attributes` _object (keys:string, values:string)_ | attributes are key/value pairs that will be added to the element when rendered. |  |  |
| `footScript` _boolean_ | footScript is a flag for script or scriptSrc that indicates if the tag should be added in the head of the page or at the foot. The default is false (add to head). To add the script to the foot of the page set footScript to true. | false |  |
| `script` _string_ | script is the text content to be added into a <script> element when rendered. |  |  |
| `scriptSrc` _string_ | scriptSrc must be an absolute URL with a protocol and host which can be used in a src attribute. |  |  |


#### TLSSpec



TLSSpec defines the desired state of TLS for a host.



_Appears in:_
- [Routing](#routing)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `secretName` _string_ | SecretName is the name of a secret that contains a TLS certificate and key. |  | MinLength: 5 <br />Required: \{\} <br /> |


#### Translation







_Appears in:_
- [KDexTranslationSpec](#kdextranslationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `keysAndValues` _object (keys:string, values:string)_ | keysAndValues is a map of key=/value pairs where the key is the identifier and the value is the translation of that key in the language specified by the lang property. |  | MinProperties: 1 <br />Required: \{\} <br /> |
| `lang` _string_ | lang is a string containing a BCP 47 language tag that identifies the set of translations.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag. |  | Required: \{\} <br /> |


