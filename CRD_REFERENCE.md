# API Reference

## Packages
- [kdex.dev/v1alpha1](#kdexdevv1alpha1)


## kdex.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group.

### Resource Types
- [MicroFrontEndApp](#microfrontendapp)
- [MicroFrontEndAppList](#microfrontendapplist)
- [MicroFrontEndHost](#microfrontendhost)
- [MicroFrontEndHostList](#microfrontendhostlist)
- [MicroFrontEndPageArchetype](#microfrontendpagearchetype)
- [MicroFrontEndPageArchetypeList](#microfrontendpagearchetypelist)
- [MicroFrontEndPageBinding](#microfrontendpagebinding)
- [MicroFrontEndPageBindingList](#microfrontendpagebindinglist)
- [MicroFrontEndPageFooter](#microfrontendpagefooter)
- [MicroFrontEndPageFooterList](#microfrontendpagefooterlist)
- [MicroFrontEndPageHeader](#microfrontendpageheader)
- [MicroFrontEndPageHeaderList](#microfrontendpageheaderlist)
- [MicroFrontEndPageNavigation](#microfrontendpagenavigation)
- [MicroFrontEndPageNavigationList](#microfrontendpagenavigationlist)
- [MicroFrontEndRenderPage](#microfrontendrenderpage)
- [MicroFrontEndRenderPageList](#microfrontendrenderpagelist)
- [MicroFrontEndTranslation](#microfrontendtranslation)
- [MicroFrontEndTranslationList](#microfrontendtranslationlist)



#### AppPolicy

_Underlying type:_ _string_

AppPolicy defines the policy for apps.

_Validation:_
- Enum: [Strict NonStrict]

_Appears in:_
- [MicroFrontEndHostSpec](#microfrontendhostspec)

| Field | Description |
| --- | --- |
| `Strict` | StrictAppPolicy means that apps may not embed JavaScript dependencies.<br /> |
| `NonStrict` | NonStrictAppPolicy means that apps may embed JavaScript dependencies.<br /> |






#### ContentEntry







_Appears in:_
- [MicroFrontEndPageBindingSpec](#microfrontendpagebindingspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `customElementName` _string_ | customElementName is the name of the MicroFrontEndApp custom element to render in the specified slot (if present in the template). |  |  |
| `appRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | appRef is a reference to the MicroFrontEndApp to include in this binding. |  |  |
| `rawHTML` _string_ | rawHTML is a raw HTML string to be rendered in the specified slot (if present in the template). |  |  |
| `slot` _string_ | slot is the name of the App slot to which this entry will be bound. If omitted, the slot used will be `main`. No more than one entry can be bound to a slot. |  |  |


#### CustomElement



CustomElement defines a custom element exposed by a micro-frontend application.



_Appears in:_
- [MicroFrontEndAppSpec](#microfrontendappspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `description` _string_ | description of the custom element. |  |  |
| `name` _string_ | name of the custom element. |  | Required: \{\} <br /> |


#### MicroFrontEndApp



MicroFrontEndApp is the Schema for the microfrontendapps API



_Appears in:_
- [MicroFrontEndAppList](#microfrontendapplist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndApp` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MicroFrontEndAppSpec](#microfrontendappspec)_ | spec defines the desired state of MicroFrontEndApp |  | Required: \{\} <br /> |


#### MicroFrontEndAppList



MicroFrontEndAppList contains a list of MicroFrontEndApp





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndAppList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndApp](#microfrontendapp) array_ |  |  |  |


#### MicroFrontEndAppSpec



MicroFrontEndAppSpec defines the desired state of MicroFrontEndApp



_Appears in:_
- [MicroFrontEndApp](#microfrontendapp)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `customElements` _[CustomElement](#customelement) array_ | customElements is a list of custom elements implemented by the micro-frontend application. |  |  |
| `packageReference` _[PackageReference](#packagereference)_ | packageReference specifies the name and version of an NPM package that contains the micro-frontend application. The package must have a package.json that contains ES modules. |  | Required: \{\} <br /> |




#### MicroFrontEndHost



MicroFrontEndHost is the Schema for the microfrontendhosts API



_Appears in:_
- [MicroFrontEndHostList](#microfrontendhostlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndHost` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MicroFrontEndHostSpec](#microfrontendhostspec)_ | spec defines the desired state of MicroFrontEndHost |  |  |


#### MicroFrontEndHostList



MicroFrontEndHostList contains a list of MicroFrontEndHost





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndHostList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndHost](#microfrontendhost) array_ |  |  |  |


#### MicroFrontEndHostSpec



MicroFrontEndHostSpec defines the desired state of MicroFrontEndHost



_Appears in:_
- [MicroFrontEndHost](#microfrontendhost)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `appPolicy` _[AppPolicy](#apppolicy)_ | AppPolicy defines the policy for apps.<br />When the strict policy is enabled, an app may not embed JavaScript dependencies.<br />Validation of the application source code will fail if dependencies are not fully externalized.<br />A Host which defines the `script` app policy must not accept apps which do not comply.<br />While a non-strict Host may accept both strict and non-strict apps. |  | Enum: [Strict NonStrict] <br />Required: \{\} <br /> |
| `baseMeta` _string_ | baseMeta is a string containing a base set of meta tags to use on every page rendered for the host. |  | MinLength: 5 <br /> |
| `defaultLang` _string_ | defaultLang is a string containing a BCP 47 language tag.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.<br />When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'. |  |  |
| `domains` _string array_ | domains are the names by which this host is addressed. The first domain listed is the preferred domain. The domains may contain wildcard prefix in the form '*.'. Longest match always wins. |  | MinItems: 1 <br />Required: \{\} <br /> |
| `organization` _string_ | Organization is the name of the Organization. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `stylesheet` _string_ | Stylesheet is the URL to the default stylesheet. |  | MinLength: 5 <br />Pattern: `^https?://` <br /> |
| `supportedLangs` _string array_ | supportedLangs is an array of strings containing BCP 47 language tags.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.<br />Render pages will be pre-rendered for each of the supported languages.<br />When not set the default will be `["en"]`. |  |  |




#### MicroFrontEndPageArchetype



MicroFrontEndPageArchetype is the Schema for the microfrontendpagearchetypes API



_Appears in:_
- [MicroFrontEndPageArchetypeList](#microfrontendpagearchetypelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageArchetype` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MicroFrontEndPageArchetypeSpec](#microfrontendpagearchetypespec)_ | spec defines the desired state of MicroFrontEndPageArchetype |  | Required: \{\} <br /> |


#### MicroFrontEndPageArchetypeList



MicroFrontEndPageArchetypeList contains a list of MicroFrontEndPageArchetype





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageArchetypeList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndPageArchetype](#microfrontendpagearchetype) array_ |  |  |  |


#### MicroFrontEndPageArchetypeSpec



MicroFrontEndPageArchetypeSpec defines the desired state of MicroFrontEndPageArchetype



_Appears in:_
- [MicroFrontEndPageArchetype](#microfrontendpagearchetype)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the structure of an HTML page. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `defaultFooterRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultFooterRef is an optional reference to a MicroFrontEndPageFooter resource. If not specified, no footer will be displayed. Use the `.Footer` property to position its content in the template. |  |  |
| `defaultHeaderRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultHeaderRef is an optional reference to a MicroFrontEndPageHeader resource. If not specified, no header will be displayed. Use the `.Header` property to position its content in the template. |  |  |
| `defaultMainNavigationRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultMainNavigationRef is an optional reference to a MicroFrontEndPageNavigation resource. If not specified, no navigation will be displayed. Use the `.Navigation["main"]` property to position its content in the template. |  |  |
| `extraNavigations` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | extraNavigations is an optional map of named navigation object references. Use `.Navigation["<name>"]` to position the named navigation's content in the template. |  |  |




#### MicroFrontEndPageBinding



MicroFrontEndPageBinding is the Schema for the microfrontendpagebindings API



_Appears in:_
- [MicroFrontEndPageBindingList](#microfrontendpagebindinglist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageBinding` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MicroFrontEndPageBindingSpec](#microfrontendpagebindingspec)_ | spec defines the desired state of MicroFrontEndPageBinding |  | Required: \{\} <br /> |


#### MicroFrontEndPageBindingList



MicroFrontEndPageBindingList contains a list of MicroFrontEndPageBinding





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageBindingList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndPageBinding](#microfrontendpagebinding) array_ |  |  |  |


#### MicroFrontEndPageBindingSpec



MicroFrontEndPageBindingSpec defines the desired state of MicroFrontEndPageBinding



_Appears in:_
- [MicroFrontEndPageBinding](#microfrontendpagebinding)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `contentEntries` _[ContentEntry](#contententry) array_ | contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or MicroFrontEndApp references. |  | MaxItems: 8 <br />MinItems: 1 <br />Required: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the MicroFrontEndHost that this binding is for. |  | Required: \{\} <br /> |
| `label` _string_ | label is the value used in menus and page titles before localization occurs (or when no translation exists for the current language). |  | Required: \{\} <br /> |
| `overrideFooterRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideFooterRef is an optional reference to a MicroFrontEndPageFooter resource. If not specified, the footer from the archetype will be used. |  |  |
| `overrideHeaderRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideHeaderRef is an optional reference to a MicroFrontEndPageHeader resource. If not specified, the header from the archetype will be used. |  |  |
| `overrideMainNavigationRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideMainNavigationRef is an optional reference to a MicroFrontEndPageNavigation resource. If not specified, the main navigation from the archetype will be used. |  |  |
| `navigationHints` _[NavigationHints](#navigationhints)_ | navigationHints are optional navigation properties that if omitted result in the page being hidden from the navigation. |  |  |
| `pageArchetypeRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | pageArchetypeRef is a reference to the MicroFrontEndPageArchetype that this binding is for. |  | Required: \{\} <br /> |
| `parentPageRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | parentPageRef is a reference to the MicroFrontEndPageBinding bellow which this page will appear in the main navigation. If not set, the page will be placed in the top level of the navigation. |  |  |
| `path` _string_ | path is the URI path at which the page will be accessible in the application server context. The final absolute path will contain this path and may be prefixed by additional context like a language identifier. |  | Required: \{\} <br /> |




#### MicroFrontEndPageFooter



MicroFrontEndPageFooter is the Schema for the microfrontendpagefooters API



_Appears in:_
- [MicroFrontEndPageFooterList](#microfrontendpagefooterlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageFooter` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MicroFrontEndPageFooterSpec](#microfrontendpagefooterspec)_ | spec defines the desired state of MicroFrontEndPageFooter |  |  |


#### MicroFrontEndPageFooterList



MicroFrontEndPageFooterList contains a list of MicroFrontEndPageFooter





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageFooterList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndPageFooter](#microfrontendpagefooter) array_ |  |  |  |


#### MicroFrontEndPageFooterSpec



MicroFrontEndPageFooterSpec defines the desired state of MicroFrontEndPageFooter



_Appears in:_
- [MicroFrontEndPageFooter](#microfrontendpagefooter)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the content of an App Server page footer section. Use the `.Footer` property to position its content in the template. |  | MinLength: 5 <br />Required: \{\} <br /> |




#### MicroFrontEndPageHeader



MicroFrontEndPageHeader is the Schema for the microfrontendpageheaders API



_Appears in:_
- [MicroFrontEndPageHeaderList](#microfrontendpageheaderlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageHeader` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MicroFrontEndPageHeaderSpec](#microfrontendpageheaderspec)_ | spec defines the desired state of MicroFrontEndPageHeader |  | Required: \{\} <br /> |


#### MicroFrontEndPageHeaderList



MicroFrontEndPageHeaderList contains a list of MicroFrontEndPageHeader





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageHeaderList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndPageHeader](#microfrontendpageheader) array_ |  |  |  |


#### MicroFrontEndPageHeaderSpec



MicroFrontEndPageHeaderSpec defines the desired state of MicroFrontEndPageHeader



_Appears in:_
- [MicroFrontEndPageHeader](#microfrontendpageheader)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the content of an App Server page header section. Use the `.Header` property to position its content in the template. |  | MinLength: 5 <br />Required: \{\} <br /> |




#### MicroFrontEndPageNavigation



MicroFrontEndPageNavigation is the Schema for the microfrontendpagenavigations API



_Appears in:_
- [MicroFrontEndPageNavigationList](#microfrontendpagenavigationlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageNavigation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MicroFrontEndPageNavigationSpec](#microfrontendpagenavigationspec)_ | spec defines the desired state of MicroFrontEndPageNavigation |  |  |


#### MicroFrontEndPageNavigationList



MicroFrontEndPageNavigationList contains a list of MicroFrontEndPageNavigation





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndPageNavigationList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndPageNavigation](#microfrontendpagenavigation) array_ |  |  |  |


#### MicroFrontEndPageNavigationSpec



MicroFrontEndPageNavigationSpec defines the desired state of MicroFrontEndPageNavigation



_Appears in:_
- [MicroFrontEndPageNavigation](#microfrontendpagenavigation)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the content of an App Server page navigation. Use the `.Navigation["<name>"]` property to position its content in the template. |  | MinLength: 5 <br />Required: \{\} <br /> |




#### MicroFrontEndRenderPage



MicroFrontEndRenderPage is the Schema for the microfrontendrenderpages API



_Appears in:_
- [MicroFrontEndRenderPageList](#microfrontendrenderpagelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndRenderPage` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MicroFrontEndRenderPageSpec](#microfrontendrenderpagespec)_ | spec defines the desired state of MicroFrontEndRenderPage |  |  |


#### MicroFrontEndRenderPageList



MicroFrontEndRenderPageList contains a list of MicroFrontEndRenderPage





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndRenderPageList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndRenderPage](#microfrontendrenderpage) array_ |  |  |  |


#### MicroFrontEndRenderPageSpec



MicroFrontEndRenderPageSpec defines the desired state of MicroFrontEndRenderPage



_Appears in:_
- [MicroFrontEndRenderPage](#microfrontendrenderpage)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the MicroFrontEndHost that this render page is for. |  | Required: \{\} <br /> |
| `navigationHints` _[NavigationHints](#navigationhints)_ | navigationHints are optional navigation properties that if omitted result in the page being hidden from the navigation. |  |  |
| `pageComponents` _[PageComponents](#pagecomponents)_ | pageComponents make up the elements of an HTML page that will be rendered by a web server. |  | Required: \{\} <br /> |
| `parentPageRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | parentPageRef is a reference to the MicroFrontEndRenderPage bellow which this page will appear in the main navigation. If not set, the page will be placed in the top level of the navigation. |  |  |
| `path` _string_ | path is the URI path at which the page will be accessible in the application server context. The final absolute path will contain this path and may be prefixed by additional context like a language identifier. |  | Required: \{\} <br /> |




#### MicroFrontEndTranslation



MicroFrontEndTranslation is the Schema for the microfrontendtranslations API



_Appears in:_
- [MicroFrontEndTranslationList](#microfrontendtranslationlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndTranslation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[MicroFrontEndTranslationSpec](#microfrontendtranslationspec)_ | spec defines the desired state of MicroFrontEndTranslation |  |  |


#### MicroFrontEndTranslationList



MicroFrontEndTranslationList contains a list of MicroFrontEndTranslation





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndTranslationList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndTranslation](#microfrontendtranslation) array_ |  |  |  |


#### MicroFrontEndTranslationSpec



MicroFrontEndTranslationSpec defines the desired state of MicroFrontEndTranslation



_Appears in:_
- [MicroFrontEndTranslation](#microfrontendtranslation)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the MicroFrontEndHost that this render page is for. |  | Required: \{\} <br /> |
| `translations` _[Translation](#translation) array_ | translations is an array of objects where each one specifies a language and a map consisting of key/value pairs. |  | MinItems: 1 <br />Required: \{\} <br /> |




#### NavigationHints







_Appears in:_
- [MicroFrontEndPageBindingSpec](#microfrontendpagebindingspec)
- [MicroFrontEndRenderPageSpec](#microfrontendrenderpagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `icon` _string_ | icon is the name of the icon to display next to the menu entry for this page. |  |  |
| `weight` _[Quantity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#quantity-resource-api)_ | weight is a property that influences the position of the page menu entry. Items at each level are sorted first by ascending weight and then ascending lexicographically. |  |  |


#### PackageReference



PackageReference specifies the name and version of an NPM package that contains the micro-frontend application.



_Appears in:_
- [MicroFrontEndAppSpec](#microfrontendappspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `secretRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | secretRef is a reference to a secret containing authentication credentials for the NPM registry that holds the package. |  |  |
| `name` _string_ | name contains a scoped npm package name. |  | Required: \{\} <br /> |
| `version` _string_ | version contains a specific npm package version. |  | Required: \{\} <br /> |


#### PageComponents







_Appears in:_
- [MicroFrontEndRenderPageSpec](#microfrontendrenderpagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `contents` _object (keys:string, values:string)_ |  |  |  |
| `footer` _string_ |  |  |  |
| `header` _string_ |  |  |  |
| `navigations` _object (keys:string, values:string)_ |  |  |  |
| `primaryTemplate` _string_ |  |  |  |
| `title` _string_ |  |  |  |


#### Translation







_Appears in:_
- [MicroFrontEndTranslationSpec](#microfrontendtranslationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `keysAndValues` _object (keys:string, values:string)_ | keysAndValues is a map of key=/value pairs where the key is the identifier and the value is the translation of that key in the language specified by the lang property. |  | MinProperties: 1 <br />Required: \{\} <br /> |
| `lang` _string_ | lang is a string containing a BCP 47 language tag that identifies the set of translations.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag. |  | Required: \{\} <br /> |


