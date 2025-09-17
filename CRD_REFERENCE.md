# API Reference

## Packages
- [kdex.dev/v1alpha1](#kdexdevv1alpha1)


## kdex.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group.

### Resource Types
- [MicroFrontEndApp](#microfrontendapp)
- [MicroFrontEndAppList](#microfrontendapplist)
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
| `status` _[MicroFrontEndAppStatus](#microfrontendappstatus)_ | status defines the observed state of MicroFrontEndApp |  |  |


#### MicroFrontEndAppList



MicroFrontEndAppList contains a list of MicroFrontEndApp





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `MicroFrontEndAppList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[MicroFrontEndApp](#microfrontendapp) array_ |  |  |  |


#### MicroFrontEndAppSource



MicroFrontEndAppSource defines the source of a micro-frontend application.



_Appears in:_
- [MicroFrontEndAppSpec](#microfrontendappspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `secretRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | secretRef is a reference to a secret containing authentication credentials for the source. |  |  |
| `url` _string_ | url of the application source. This can be a Git repository, an archive, or an OCI artifact. |  | Required: \{\} <br /> |


#### MicroFrontEndAppSpec



MicroFrontEndAppSpec defines the desired state of MicroFrontEndApp



_Appears in:_
- [MicroFrontEndApp](#microfrontendapp)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `customElements` _[CustomElement](#customelement) array_ | customElements is a list of custom elements implemented by the micro-frontend application. |  |  |
| `source` _[MicroFrontEndAppSource](#microfrontendappsource)_ | source configures the location of the source code of the micro-frontend application. The source code must contain a valid package.json that produces ES modules. Based on App Server configuration embedded dependencies may not be allowed. In this case dependencies must be externalized otherwise the app CR will not validate. |  | Required: \{\} <br /> |


#### MicroFrontEndAppStatus



MicroFrontEndAppStatus defines the observed state of MicroFrontEndApp.



_Appears in:_
- [MicroFrontEndApp](#microfrontendapp)



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
| `status` _[MicroFrontEndPageArchetypeStatus](#microfrontendpagearchetypestatus)_ | status defines the observed state of MicroFrontEndPageArchetype |  |  |


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
| `content` _string_ | content is a go string template that defines the structure of an App Server page. The template accesses `.Values` properties to render its contents. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `defaultFooterRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultFooterRef is an optional reference to a MicroFrontEndPageFooter resource. If not specified, no footer will be displayed. |  |  |
| `defaultHeaderRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultHeaderRef is an optional reference to a MicroFrontEndPageHeader resource. If not specified, no header will be displayed. |  |  |
| `defaultMainNavigationRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | defaultMainNavigationRef is an optional reference to a MicroFrontEndPageNavigation resource referenced as `\{\{ .Values.navigation["main"] \}\}`. If not specified, no navigation will be displayed. |  |  |
| `extraNavigations` _object (keys:string, values:[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core))_ | extraNavigations is an optional map of named navigation object references that will be available in page templates as `\{\{ .Values.navigation["name"] \}\}`. |  |  |


#### MicroFrontEndPageArchetypeStatus



MicroFrontEndPageArchetypeStatus defines the observed state of MicroFrontEndPageArchetype.



_Appears in:_
- [MicroFrontEndPageArchetype](#microfrontendpagearchetype)



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
| `status` _[MicroFrontEndPageBindingStatus](#microfrontendpagebindingstatus)_ | status defines the observed state of MicroFrontEndPageBinding |  |  |


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
| `label` _string_ | label is the value used in menus and page titles before localization occurs (or when no translation exists for the current language). |  | Required: \{\} <br /> |
| `contentEntries` _[ContentEntry](#contententry) array_ | contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or MicroFrontEndApp references. |  | MaxItems: 8 <br />MinItems: 1 <br />Required: \{\} <br /> |
| `pageArchetypeRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | pageArchetypeRef is a reference to the MicroFrontEndPageArchetype that this binding is for. |  | Required: \{\} <br /> |
| `overrideFooterRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideFooterRef is an optional reference to a MicroFrontEndPageFooter resource. If not specified, the footer from the archetype will be used. |  |  |
| `overrideHeaderRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideHeaderRef is an optional reference to a MicroFrontEndPageHeader resource. If not specified, the header from the archetype will be used. |  |  |
| `overrideMainNavigationRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | overrideMainNavigationRef is an optional reference to a MicroFrontEndPageNavigation resource referenced as `\{\{ .Values.navigation["main"] \}\}. If not specified, the main navigation from the archetype will be used. |  |  |
| `parent` _string_ | parent specifies the menu entry that is the parent under which the menu entry for this page will be added in the main navigation. A hierarchical path using slashes is supported. |  |  |
| `path` _string_ | path is the URI path at which the page will be accessible in the application server context. The final absolute path will contain this path and may be prefixed by additional context like a language identifier. |  | Required: \{\} <br /> |
| `weight` _[Quantity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#quantity-resource-api)_ | weight is a property that influences the position of the page menu entry. Items are sorted first by ascending weight and then ascending lexicographically. |  |  |


#### MicroFrontEndPageBindingStatus



MicroFrontEndPageBindingStatus defines the observed state of MicroFrontEndPageBinding.



_Appears in:_
- [MicroFrontEndPageBinding](#microfrontendpagebinding)



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
| `status` _[MicroFrontEndPageFooterStatus](#microfrontendpagefooterstatus)_ | status defines the observed state of MicroFrontEndPageFooter |  |  |


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
| `content` _string_ | content is a go string template that defines the content of an App Server page footer section. The template accesses `.Values` properties to render its contents. |  | MinLength: 5 <br />Required: \{\} <br /> |


#### MicroFrontEndPageFooterStatus



MicroFrontEndPageFooterStatus defines the observed state of MicroFrontEndPageFooter.



_Appears in:_
- [MicroFrontEndPageFooter](#microfrontendpagefooter)



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
| `status` _[MicroFrontEndPageHeaderStatus](#microfrontendpageheaderstatus)_ | status defines the observed state of MicroFrontEndPageHeader |  |  |


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
| `content` _string_ | content is a go string template that defines the content of an App Server page header section. The template accesses `.Values` properties to render its contents. |  | MinLength: 5 <br />Required: \{\} <br /> |


#### MicroFrontEndPageHeaderStatus



MicroFrontEndPageHeaderStatus defines the observed state of MicroFrontEndPageHeader.



_Appears in:_
- [MicroFrontEndPageHeader](#microfrontendpageheader)



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
| `status` _[MicroFrontEndPageNavigationStatus](#microfrontendpagenavigationstatus)_ | status defines the observed state of MicroFrontEndPageNavigation |  |  |


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
| `content` _string_ | content is a go string template that defines the content of an App Server page navigation. The template accesses `.Values` properties to render its contents. |  | MinLength: 5 <br />Required: \{\} <br /> |


#### MicroFrontEndPageNavigationStatus



MicroFrontEndPageNavigationStatus defines the observed state of MicroFrontEndPageNavigation.



_Appears in:_
- [MicroFrontEndPageNavigation](#microfrontendpagenavigation)



