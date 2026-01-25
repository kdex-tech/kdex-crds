# API Reference

## Packages
- [kdex.dev/v1alpha1](#kdexdevv1alpha1)


## kdex.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group.

### Resource Types
- [KDexApp](#kdexapp)
- [KDexAppList](#kdexapplist)
- [KDexClusterApp](#kdexclusterapp)
- [KDexClusterAppList](#kdexclusterapplist)
- [KDexClusterPageArchetype](#kdexclusterpagearchetype)
- [KDexClusterPageArchetypeList](#kdexclusterpagearchetypelist)
- [KDexClusterPageFooter](#kdexclusterpagefooter)
- [KDexClusterPageFooterList](#kdexclusterpagefooterlist)
- [KDexClusterPageHeader](#kdexclusterpageheader)
- [KDexClusterPageHeaderList](#kdexclusterpageheaderlist)
- [KDexClusterPageNavigation](#kdexclusterpagenavigation)
- [KDexClusterPageNavigationList](#kdexclusterpagenavigationlist)
- [KDexClusterScriptLibrary](#kdexclusterscriptlibrary)
- [KDexClusterScriptLibraryList](#kdexclusterscriptlibrarylist)
- [KDexClusterTheme](#kdexclustertheme)
- [KDexClusterThemeList](#kdexclusterthemelist)
- [KDexClusterTranslation](#kdexclustertranslation)
- [KDexClusterTranslationList](#kdexclustertranslationlist)
- [KDexClusterUtilityPage](#kdexclusterutilitypage)
- [KDexClusterUtilityPageList](#kdexclusterutilitypagelist)
- [KDexFunction](#kdexfunction)
- [KDexFunctionList](#kdexfunctionlist)
- [KDexHost](#kdexhost)
- [KDexHostList](#kdexhostlist)
- [KDexInternalHost](#kdexinternalhost)
- [KDexInternalHostList](#kdexinternalhostlist)
- [KDexInternalPackageReferences](#kdexinternalpackagereferences)
- [KDexInternalPackageReferencesList](#kdexinternalpackagereferenceslist)
- [KDexInternalTranslation](#kdexinternaltranslation)
- [KDexInternalTranslationList](#kdexinternaltranslationlist)
- [KDexInternalUtilityPage](#kdexinternalutilitypage)
- [KDexInternalUtilityPageList](#kdexinternalutilitypagelist)
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
- [KDexScope](#kdexscope)
- [KDexScopeBinding](#kdexscopebinding)
- [KDexScopeBindingList](#kdexscopebindinglist)
- [KDexScopeList](#kdexscopelist)
- [KDexScriptLibrary](#kdexscriptlibrary)
- [KDexScriptLibraryList](#kdexscriptlibrarylist)
- [KDexTheme](#kdextheme)
- [KDexThemeList](#kdexthemelist)
- [KDexTranslation](#kdextranslation)
- [KDexTranslationList](#kdextranslationlist)
- [KDexUtilityPage](#kdexutilitypage)
- [KDexUtilityPageList](#kdexutilitypagelist)



#### API







_Appears in:_
- [KDexFunctionSpec](#kdexfunctionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `basePath` _string_ | basePath is the base URL path for the function. It must match the regex ^/\w+/\w+ (e.g., /v1/users). |  | Pattern: `^/\w+/\w+` <br />Required: \{\} <br /> |
| `paths` _object (keys:string, values:[PathItem](#pathitem))_ | paths is a map of paths that exist below the basePath. All keys of the map must be paths prefixed by .spec.api.basePath. |  | MaxProperties: 16 <br />MinProperties: 1 <br />Required: \{\} <br /> |
| `schemas` _object (keys:string, values:[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg))_ |  |  | MaxProperties: 6 <br />Optional: \{\} <br /> |


#### Asset







_Appears in:_
- [Assets](#assets)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `attributes` _object (keys:string, values:string)_ | attributes are key/value pairs that will be added to the element as attributes when rendered. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |
| `linkHref` _string_ | linkHref is the content of a `<link>` href attribute. The URL may be absolute with protocol and host or it must be prefixed by the IngressPath of the Backend. |  | Optional: \{\} <br /> |
| `metaId` _string_ | metaId is required just for semantics of CRD field validation. |  | Optional: \{\} <br /> |
| `style` _string_ | style is the text content to be added into a `<style>` element when rendered. |  | Optional: \{\} <br /> |


#### Assets

_Underlying type:_ _[Asset](#asset)_



_Validation:_
- MaxItems: 32

_Appears in:_
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)
- [KDexThemeSpec](#kdexthemespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `attributes` _object (keys:string, values:string)_ | attributes are key/value pairs that will be added to the element as attributes when rendered. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |
| `linkHref` _string_ | linkHref is the content of a `<link>` href attribute. The URL may be absolute with protocol and host or it must be prefixed by the IngressPath of the Backend. |  | Optional: \{\} <br /> |
| `metaId` _string_ | metaId is required just for semantics of CRD field validation. |  | Optional: \{\} <br /> |
| `style` _string_ | style is the text content to be added into a `<style>` element when rendered. |  | Optional: \{\} <br /> |


#### Auth







_Appears in:_
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `jwt` _[JWT](#jwt)_ | jwt is the configuation for JWT token support. |  | Optional: \{\} <br /> |
| `oidcProvider` _[OIDCProvider](#oidcprovider)_ | oidcProvider is the configuration for an optional OIDC provider. |  | Optional: \{\} <br /> |


#### Backend



Backend defines a deployment for serving resources specific to the refer.



_Appears in:_
- [KDexAppSpec](#kdexappspec)
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)
- [KDexScriptLibrarySpec](#kdexscriptlibraryspec)
- [KDexThemeSpec](#kdexthemespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | imagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.<br />More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod |  | Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/_' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/_.+` <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |
| `serverImage` _string_ | serverImage is the name of Backend image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `serverImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the Backend server image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `staticImage` _string_ | staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `staticImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |








#### ContactInfo



ContactInfo defines contact details.



_Appears in:_
- [KDexFunctionMetadata](#kdexfunctionmetadata)
- [KDexPageBindingSpec](#kdexpagebindingspec)
- [Metadata](#metadata)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name of the contact. |  | Optional: \{\} <br /> |
| `email` _string_ | Email of the contact. |  | Optional: \{\} <br /> |


#### ContentEntry







_Appears in:_
- [KDexInternalUtilityPageSpec](#kdexinternalutilitypagespec)
- [KDexPageBindingSpec](#kdexpagebindingspec)
- [KDexUtilityPageSpec](#kdexutilitypagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `slot` _string_ | slot is the unique name to which this entry will be bound. |  | Required: \{\} <br /> |
| `appRef` _[KDexObjectReference](#kdexobjectreference)_ | appRef is a reference to the KDexApp to include in this binding. |  | Optional: \{\} <br /> |
| `customElementName` _string_ | customElementName is the name of the KDexApp custom element to render in the specified slot (if present in the template). |  | Optional: \{\} <br /> |
| `attributes` _object (keys:string, values:string)_ | attributes are key/value pairs that will be added to the custom element as attributes when rendered. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |
| `rawHTML` _string_ | rawHTML is a raw HTML string to be rendered in the specified slot (if present in the template). |  | Optional: \{\} <br /> |


#### ContentEntryApp







_Appears in:_
- [ContentEntry](#contententry)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `appRef` _[KDexObjectReference](#kdexobjectreference)_ | appRef is a reference to the KDexApp to include in this binding. |  | Optional: \{\} <br /> |
| `customElementName` _string_ | customElementName is the name of the KDexApp custom element to render in the specified slot (if present in the template). |  | Optional: \{\} <br /> |
| `attributes` _object (keys:string, values:string)_ | attributes are key/value pairs that will be added to the custom element as attributes when rendered. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |


#### ContentEntryStatic







_Appears in:_
- [ContentEntry](#contententry)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `rawHTML` _string_ | rawHTML is a raw HTML string to be rendered in the specified slot (if present in the template). |  | Optional: \{\} <br /> |


#### CustomElement



CustomElement defines a custom element exposed by a micro-frontend application.



_Appears in:_
- [KDexAppSpec](#kdexappspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | name of the custom element. |  | Required: \{\} <br /> |
| `description` _string_ | description of the custom element. |  | Optional: \{\} <br /> |


#### JWT







_Appears in:_
- [Auth](#auth)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `activeKey` _string_ | activeKey contains the name of the secret that holds the currently active key. This can be omitted when there is only a single key specified. |  | Optional: \{\} <br /> |
| `jwtKeysSecrets` _[LocalSecretWithKeyReference](#localsecretwithkeyreference) array_ | jwtKeysSecrets is an optional list of references to secrets in the same namespace that hold private PEM encoded signing keys. |  | Optional: \{\} <br /> |
| `tokenTTL` _string_ | tokenTTL is the length of time for which the token is valid | 1h | Optional: \{\} <br /> |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
- [KDexClusterApp](#kdexclusterapp)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `customElements` _[CustomElement](#customelement) array_ | customElements is a list of custom elements implemented by the micro-frontend application. |  | MaxItems: 32 <br />MinItems: 1 <br /> |
| `packageReference` _[PackageReference](#packagereference)_ | packageReference specifies the name and version of an NPM package that contains the script. The package.json must describe an ES module. |  | Required: \{\} <br /> |
| `scripts` _[ScriptDef](#scriptdef) array_ | scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents. |  | MaxItems: 8 <br />Optional: \{\} <br /> |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | imagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.<br />More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod |  | Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/_' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/_.+` <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |
| `serverImage` _string_ | serverImage is the name of Backend image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `serverImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the Backend server image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `staticImage` _string_ | staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `staticImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |


#### KDexClusterApp



KDexClusterApp is the Schema for the kdexclusterapps API



_Appears in:_
- [KDexClusterAppList](#kdexclusterapplist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterApp` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexAppSpec](#kdexappspec)_ | spec defines the desired state of KDexClusterApp |  |  |


#### KDexClusterAppList



KDexClusterAppList contains a list of KDexClusterApp





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterAppList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterApp](#kdexclusterapp) array_ |  |  |  |


#### KDexClusterPageArchetype



KDexClusterPageArchetype is the Schema for the kdexclusterpagearchetypes API



_Appears in:_
- [KDexClusterPageArchetypeList](#kdexclusterpagearchetypelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterPageArchetype` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexPageArchetypeSpec](#kdexpagearchetypespec)_ | spec defines the desired state of KDexClusterPageArchetype |  |  |


#### KDexClusterPageArchetypeList



KDexClusterPageArchetypeList contains a list of KDexClusterPageArchetype





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterPageArchetypeList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterPageArchetype](#kdexclusterpagearchetype) array_ |  |  |  |


#### KDexClusterPageFooter



KDexClusterPageFooter is the Schema for the kdexclusterpagefooters API



_Appears in:_
- [KDexClusterPageFooterList](#kdexclusterpagefooterlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterPageFooter` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexPageFooterSpec](#kdexpagefooterspec)_ | spec defines the desired state of KDexClusterPageFooter |  |  |


#### KDexClusterPageFooterList



KDexClusterPageFooterList contains a list of KDexClusterPageFooter





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterPageFooterList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterPageFooter](#kdexclusterpagefooter) array_ |  |  |  |


#### KDexClusterPageHeader



KDexClusterPageHeader is the Schema for the kdexclusterpageheaders API



_Appears in:_
- [KDexClusterPageHeaderList](#kdexclusterpageheaderlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterPageHeader` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexPageHeaderSpec](#kdexpageheaderspec)_ | spec defines the desired state of KDexClusterPageHeader |  |  |


#### KDexClusterPageHeaderList



KDexClusterPageHeaderList contains a list of KDexClusterPageHeader





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterPageHeaderList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterPageHeader](#kdexclusterpageheader) array_ |  |  |  |


#### KDexClusterPageNavigation



KDexClusterPageNavigation is the Schema for the kdexclusterpagenavigations API



_Appears in:_
- [KDexClusterPageNavigationList](#kdexclusterpagenavigationlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterPageNavigation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexPageNavigationSpec](#kdexpagenavigationspec)_ | spec defines the desired state of KDexClusterPageNavigation |  |  |


#### KDexClusterPageNavigationList



KDexClusterPageNavigationList contains a list of KDexClusterPageNavigation





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterPageNavigationList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterPageNavigation](#kdexclusterpagenavigation) array_ |  |  |  |


#### KDexClusterScriptLibrary



KDexClusterScriptLibrary is the Schema for the kdexclusterscriptlibraries API



_Appears in:_
- [KDexClusterScriptLibraryList](#kdexclusterscriptlibrarylist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterScriptLibrary` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexScriptLibrarySpec](#kdexscriptlibraryspec)_ | spec defines the desired state of KDexClusterScriptLibrary |  |  |


#### KDexClusterScriptLibraryList



KDexClusterScriptLibraryList contains a list of KDexClusterScriptLibrary





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterScriptLibraryList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterScriptLibrary](#kdexclusterscriptlibrary) array_ |  |  |  |


#### KDexClusterTheme



KDexClusterTheme is the Schema for the kdexclusterthemes API



_Appears in:_
- [KDexClusterThemeList](#kdexclusterthemelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterTheme` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexThemeSpec](#kdexthemespec)_ | spec defines the desired state of KDexClusterTheme |  |  |


#### KDexClusterThemeList



KDexClusterThemeList contains a list of KDexClusterTheme





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterThemeList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterTheme](#kdexclustertheme) array_ |  |  |  |


#### KDexClusterTranslation



KDexClusterTranslation is the Schema for the kdexclustertranslations API

KDexClusterTranslations allow resources to be internationalized by making translations available in as many languages
as necessary.



_Appears in:_
- [KDexClusterTranslationList](#kdexclustertranslationlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterTranslation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexTranslationSpec](#kdextranslationspec)_ | spec defines the desired state of KDexTranslation |  | Required: \{\} <br /> |


#### KDexClusterTranslationList



KDexClusterTranslationList contains a list of KDexClusterTranslation





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterTranslationList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterTranslation](#kdexclustertranslation) array_ |  |  |  |


#### KDexClusterUtilityPage



KDexClusterUtilityPage is the Schema for the kdexclusterutilitypages API

A KDexClusterUtilityPage is a cluster-scoped version of KDexUtilityPage. It allows defining default utility pages
that can be used across multiple hosts if not overridden.



_Appears in:_
- [KDexClusterUtilityPageList](#kdexclusterutilitypagelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterUtilityPage` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexUtilityPageSpec](#kdexutilitypagespec)_ | spec defines the desired state of KDexClusterUtilityPage |  | Required: \{\} <br /> |


#### KDexClusterUtilityPageList



KDexClusterUtilityPageList contains a list of KDexClusterUtilityPage





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterUtilityPageList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterUtilityPage](#kdexclusterutilitypage) array_ |  |  |  |


#### KDexFunction



KDexFunction is the Schema for the kdexfunctions API.

KDexFunction is a facility to express a particularly concise unit of logic that scales in isolation.
Ideally these are utilized via a FaaS layer, but for simplicity, some scenarios are modeled by the
Backend type using containers (native Kubernetes workloads).



_Appears in:_
- [KDexFunctionList](#kdexfunctionlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexFunction` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexFunctionSpec](#kdexfunctionspec)_ | spec defines the desired state of KDexFunction |  | Required: \{\} <br /> |


#### KDexFunctionExec



KDexFunctionExec defines the FaaS execution environment.



_Appears in:_
- [KDexFunctionSpec](#kdexfunctionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `codePackage` _string_ | CodePackage is a reference to the compiled code artifact or source package. |  | Optional: \{\} <br /> |
| `entrypoint` _string_ | Entrypoint is the specific function handler/method to execute. |  | Optional: \{\} <br /> |
| `environment` _string_ | Environment is the FaaS environment name (e.g., go-env, python-env). |  | Optional: \{\} <br /> |
| `language` _string_ | Language is the programming language of the function (e.g., go, python, nodejs). |  | Optional: \{\} <br /> |
| `scaling` _[ScalingConfig](#scalingconfig)_ | Scaling allows configuration for min/max replicas and autoscaler type. |  | Optional: \{\} <br /> |


#### KDexFunctionList



KDexFunctionList contains a list of KDexFunction





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexFunctionList` | | |
| `items` _[KDexFunction](#kdexfunction) array_ |  |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |


#### KDexFunctionMetadata



KDexFunctionMetadata defines the metadata for the function.



_Appears in:_
- [KDexFunctionSpec](#kdexfunctionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `autoGenerated` _boolean_ | AutoGenerated indicates if this CR was created by the "Sniffer" server. |  | Optional: \{\} <br /> |
| `tags` _string array_ | Tags are used for grouping and searching functions. |  | Optional: \{\} <br /> |
| `contact` _[ContactInfo](#contactinfo)_ | Contact provides contact information for the function's owner. |  | Optional: \{\} <br /> |
| `sourceImage` _string_ | SourceImage is the OCI artifact reference where the stub code was pushed. |  | Optional: \{\} <br /> |


#### KDexFunctionSpec



KDexFunctionSpec defines the desired state of KDexFunction



_Appears in:_
- [KDexFunction](#kdexfunction)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `api` _[API](#api)_ | API defines the OpenAPI contract for the function.<br />See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#path-item-object<br />See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#schema-object<br />The supported fields from 'path item object' are: summary, description, get, put, post, delete, options, head, patch, trace, parameters, and responses.<br />The field 'schemas' of type map[string]schema whose values are defined by 'schema object' is supported and can be referenced throughout operation definitions. References must be in the form "#/components/schemas/<name>". |  | Required: \{\} <br /> |
| `function` _[KDexFunctionExec](#kdexfunctionexec)_ | Function defines the FaaS execution details. |  | Optional: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this translation belongs to. |  | Required: \{\} <br /> |
| `metadata` _[KDexFunctionMetadata](#kdexfunctionmetadata)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |


#### KDexFunctionState

_Underlying type:_ _string_

KDexFunctionState reflects the current state of a KDexFunction.

_Validation:_
- Enum: [Building Pending Ready StubGenerated]

_Appears in:_
- [KDexFunctionStatus](#kdexfunctionstatus)

| Field | Description |
| --- | --- |
| `Building` | KDexFunctionStateBuilding indicates the function is being built.<br /> |
| `Pending` | KDexFunctionStatePending indicates the function is pending action.<br /> |
| `Ready` | KDexFunctionStateReady indicates the function is ready for invocation.<br /> |
| `StubGenerated` | KDexFunctionStateStubGenerated indicates the function stub has been generated.<br /> |




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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexHostSpec](#kdexhostspec)_ | spec defines the desired state of KDexHost |  | Required: \{\} <br /> |


#### KDexHostList



KDexHostList contains a list of KDexHost





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexHostList` | | |
| `items` _[KDexHost](#kdexhost) array_ | items contains a list of KDexHost |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |


#### KDexHostSpec



KDexHostSpec defines the desired state of KDexHost



_Appears in:_
- [KDexHost](#kdexhost)
- [KDexInternalHostSpec](#kdexinternalhostspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `assets` _[Assets](#assets)_ | assets is a set of elements that define a host specific HTML instructions (e.g. favicon, site logo, charset). |  | MaxItems: 32 <br /> |
| `auth` _[Auth](#auth)_ | auth holds the host's authentication configuration. |  | Optional: \{\} <br /> |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | imagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.<br />More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod |  | Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/_' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/_.+` <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |
| `serverImage` _string_ | serverImage is the name of Backend image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `serverImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the Backend server image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `staticImage` _string_ | staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `staticImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `brandName` _string_ | brandName is the name used when rendering pages belonging to the host. For example, it may be used as alt text for the logo displayed in the page header. |  | Required: \{\} <br /> |
| `defaultLang` _string_ | defaultLang is a string containing a BCP 47 language tag.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.<br />When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'. |  | Optional: \{\} <br /> |
| `devMode` _boolean_ | devMode is a boolean that enables development features like the Request Sniffer. |  | Optional: \{\} <br /> |
| `faviconSVGTemplate` _string_ | faviconSVGTemplate contains SVG code marked up with go string template to which will be passed the render.TemplateData holding other host details. The rendered output will be cached and served at "/favicon.ico" as "image/svg+xml". |  | Optional: \{\} <br /> |
| `modulePolicy` _[ModulePolicy](#modulepolicy)_ | modulePolicy defines the policy for JavaScript references in KDexApp, KDexTheme and KDexScriptLibrary resources. When not specified the policy is Strict<br />A Host must not accept JavaScript references which do not comply with the specified policy. | Strict | Enum: [ExternalDependencies Loose ModulesRequired Strict] <br />Optional: \{\} <br /> |
| `organization` _string_ | organization is the name of the Organization to which the host belongs. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `routing` _[Routing](#routing)_ | routing defines the desired routing configuration for the host. |  | Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |
| `security` _[SecurityRequirement](#securityrequirement)_ | Optional top level security requirements. |  |  |
| `themeRef` _[KDexObjectReference](#kdexobjectreference)_ | themeRef is a reference to the theme that should apply to all pages bound to this host. |  | Optional: \{\} <br /> |
| `translationRefs` _[KDexObjectReference](#kdexobjectreference) array_ | translationRefs is an array of references to KDexTranslation or KDexClusterTranslation resources that define the translations that should apply to this host. |  | Optional: \{\} <br /> |
| `utilityPages` _[UtilityPages](#utilitypages)_ | utilityPages defines the utility pages (announcement, error, login) for the host. |  | Optional: \{\} <br /> |


#### KDexInternalHost



KDexInternalHost is the Schema for the kdexinternalhosts API

A KDexInternalHost is the resource used to instantiate and manage a unique controller focused on a single KDexHost
resource. This focused controller serves to aggregate the host specific resources, primarily KDexPageBindings but
also as the main web server handling page rendering and page serving. In order to isolate the resources consumed by
those operations from other hosts a unique controller is necessary. This resource is internally generated and managed
and not meant for end users.



_Appears in:_
- [KDexInternalHostList](#kdexinternalhostlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexInternalHost` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexInternalHostSpec](#kdexinternalhostspec)_ | spec defines the desired state of KDexInternalHost |  | Required: \{\} <br /> |


#### KDexInternalHostList



KDexInternalHostList contains a list of KDexInternalHost





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexInternalHostList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexInternalHost](#kdexinternalhost) array_ |  |  |  |


#### KDexInternalHostSpec



KDexInternalHostSpec defines the desired state of KDexInternalHost



_Appears in:_
- [KDexInternalHost](#kdexinternalhost)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `assets` _[Assets](#assets)_ | assets is a set of elements that define a host specific HTML instructions (e.g. favicon, site logo, charset). |  | MaxItems: 32 <br /> |
| `auth` _[Auth](#auth)_ | auth holds the host's authentication configuration. |  | Optional: \{\} <br /> |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | imagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.<br />More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod |  | Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/_' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/_.+` <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |
| `serverImage` _string_ | serverImage is the name of Backend image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `serverImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the Backend server image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `staticImage` _string_ | staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `staticImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `brandName` _string_ | brandName is the name used when rendering pages belonging to the host. For example, it may be used as alt text for the logo displayed in the page header. |  | Required: \{\} <br /> |
| `defaultLang` _string_ | defaultLang is a string containing a BCP 47 language tag.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.<br />When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'. |  | Optional: \{\} <br /> |
| `devMode` _boolean_ | devMode is a boolean that enables development features like the Request Sniffer. |  | Optional: \{\} <br /> |
| `faviconSVGTemplate` _string_ | faviconSVGTemplate contains SVG code marked up with go string template to which will be passed the render.TemplateData holding other host details. The rendered output will be cached and served at "/favicon.ico" as "image/svg+xml". |  | Optional: \{\} <br /> |
| `modulePolicy` _[ModulePolicy](#modulepolicy)_ | modulePolicy defines the policy for JavaScript references in KDexApp, KDexTheme and KDexScriptLibrary resources. When not specified the policy is Strict<br />A Host must not accept JavaScript references which do not comply with the specified policy. | Strict | Enum: [ExternalDependencies Loose ModulesRequired Strict] <br />Optional: \{\} <br /> |
| `organization` _string_ | organization is the name of the Organization to which the host belongs. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `routing` _[Routing](#routing)_ | routing defines the desired routing configuration for the host. |  | Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |
| `security` _[SecurityRequirement](#securityrequirement)_ | Optional top level security requirements. |  |  |
| `themeRef` _[KDexObjectReference](#kdexobjectreference)_ | themeRef is a reference to the theme that should apply to all pages bound to this host. |  | Optional: \{\} <br /> |
| `translationRefs` _[KDexObjectReference](#kdexobjectreference) array_ | translationRefs is an array of references to KDexTranslation or KDexClusterTranslation resources that define the translations that should apply to this host. |  | Optional: \{\} <br /> |
| `utilityPages` _[UtilityPages](#utilitypages)_ | utilityPages defines the utility pages (announcement, error, login) for the host. |  | Optional: \{\} <br /> |
| `announcementRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | announcementRef is a reference to the KDexInternalUtilityPage that provides the announcement page. |  | Optional: \{\} <br /> |
| `errorRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | errorRef is a reference to the KDexInternalUtilityPage that provides the error page. |  | Optional: \{\} <br /> |
| `loginRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | loginRef is a reference to the KDexInternalUtilityPage that provides the login page. |  | Optional: \{\} <br /> |
| `requiredBackends` _[KDexObjectReference](#kdexobjectreference) array_ | requiredBackends is a set of references to KDexApp or KDexScriptLibrary resources that specify a backend. |  | Optional: \{\} <br /> |
| `internalTranslationRefs` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | internalTranslationRefs is a set of references to KDexInternalTranslation resources that specify a translation. |  | Optional: \{\} <br /> |


#### KDexInternalPackageReferences



KDexInternalPackageReferences is the Schema for the kdexinternalpackagereferences API

KDexInternalPackageReferences is the resource used to collect and drive the build and packaging of the complete set of npm
modules referenced by all the resources associated with a given KDexHost. This resource is internally generated and
managed and not meant for end users.



_Appears in:_
- [KDexInternalPackageReferencesList](#kdexinternalpackagereferenceslist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexInternalPackageReferences` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexInternalPackageReferencesSpec](#kdexinternalpackagereferencesspec)_ | spec defines the desired state of KDexInternalPackageReferences |  |  |


#### KDexInternalPackageReferencesList



KDexInternalPackageReferencesList contains a list of KDexInternalPackageReferences





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexInternalPackageReferencesList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexInternalPackageReferences](#kdexinternalpackagereferences) array_ |  |  |  |


#### KDexInternalPackageReferencesSpec



KDexInternalPackageReferencesSpec defines the desired state of KDexInternalPackageReferences



_Appears in:_
- [KDexInternalPackageReferences](#kdexinternalpackagereferences)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `packageReferences` _[PackageReference](#packagereference) array_ |  |  | MinItems: 1 <br /> |


#### KDexInternalTranslation



KDexInternalTranslation is the Schema for the kdexinternaltranslations API



_Appears in:_
- [KDexInternalTranslationList](#kdexinternaltranslationlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexInternalTranslation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexInternalTranslationSpec](#kdexinternaltranslationspec)_ | spec defines the desired state of KDexInternalTranslation |  | Required: \{\} <br /> |


#### KDexInternalTranslationList



KDexInternalTranslationList contains a list of KDexInternalTranslation





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexInternalTranslationList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexInternalTranslation](#kdexinternaltranslation) array_ |  |  |  |


#### KDexInternalTranslationSpec







_Appears in:_
- [KDexInternalTranslation](#kdexinternaltranslation)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `translations` _[Translation](#translation) array_ | translations is an array of objects where each one specifies a language (lang) and a map (keysAndValues) consisting of key/value pairs. If the lang property is not unique in the array and its keysAndValues map contains the same keys, the last one takes precedence. |  | MaxItems: 32 <br />MinItems: 1 <br />Required: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexInternalHost that this translation belongs to. |  | Required: \{\} <br /> |


#### KDexInternalUtilityPage



KDexInternalUtilityPage is the Schema for the kdexinternalutilitypages API

A KDexInternalUtilityPage is an internal resource used to instantiate access to a utility page for a specific host.
It is created by the KDex Host controller based on either a specific KDexUtilityPage/KDexClusterUtilityPage reference
or a default system-wide utility page.



_Appears in:_
- [KDexInternalUtilityPageList](#kdexinternalutilitypagelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexInternalUtilityPage` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexInternalUtilityPageSpec](#kdexinternalutilitypagespec)_ | spec defines the desired state of KDexInternalUtilityPage |  | Required: \{\} <br /> |


#### KDexInternalUtilityPageList



KDexInternalUtilityPageList contains a list of KDexInternalUtilityPage





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexInternalUtilityPageList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexInternalUtilityPage](#kdexinternalutilitypage) array_ |  |  |  |


#### KDexInternalUtilityPageSpec



KDexInternalUtilityPageSpec defines the desired state of KDexInternalUtilityPage



_Appears in:_
- [KDexInternalUtilityPage](#kdexinternalutilitypage)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[KDexUtilityPageType](#kdexutilitypagetype)_ | type indicates the purpose of this utility page. |  | Enum: [Announcement Error Login] <br />Required: \{\} <br /> |
| `contentEntries` _[ContentEntry](#contententry) array_ | contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or KDexApp references. |  | MaxItems: 8 <br />MinItems: 1 <br />Required: \{\} <br /> |
| `overrideFooterRef` _[KDexObjectReference](#kdexobjectreference)_ | overrideFooterRef is an optional reference to a KDexPageFooter resource. If not specified, the footer from the archetype will be used. |  | Optional: \{\} <br /> |
| `overrideHeaderRef` _[KDexObjectReference](#kdexobjectreference)_ | overrideHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, the header from the archetype will be used. |  | Optional: \{\} <br /> |
| `overrideNavigationRefs` _object (keys:string, values:[KDexObjectReference](#kdexobjectreference))_ | overrideNavigationRefs is an optional map of keyed navigation object references. When not empty, the 'main' key must be specified. These navigations will be merged with the navigations from the archetype. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |
| `pageArchetypeRef` _[KDexObjectReference](#kdexobjectreference)_ | pageArchetypeRef is a reference to the KDexPageArchetype that this binding is for. |  | Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexInternalHost that this utility page belongs to. |  | Required: \{\} <br /> |




#### KDexObjectReference







_Appears in:_
- [ContentEntry](#contententry)
- [ContentEntryApp](#contententryapp)
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)
- [KDexInternalUtilityPageSpec](#kdexinternalutilitypagespec)
- [KDexPageArchetypeSpec](#kdexpagearchetypespec)
- [KDexPageBindingSpec](#kdexpagebindingspec)
- [KDexPageFooterSpec](#kdexpagefooterspec)
- [KDexPageHeaderSpec](#kdexpageheaderspec)
- [KDexPageNavigationSpec](#kdexpagenavigationspec)
- [KDexThemeSpec](#kdexthemespec)
- [KDexUtilityPageSpec](#kdexutilitypagespec)
- [UtilityPages](#utilitypages)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name of the referent. |  | Required: \{\} <br /> |
| `kind` _string_ | Kind is the type of resource being referenced |  | Required: \{\} <br /> |
| `namespace` _string_ | Namespace, if set, causes the lookup for the namespace scoped Kind of the referent to use the specified<br />namespace. If not set, the namespace of the resource will be used to lookup the namespace scoped Kind of the<br />referent.<br />If the referring resource is cluster scoped, this field is ignored.<br />Defaulted to nil. |  | Optional: \{\} <br /> |


#### KDexObjectStatus







_Appears in:_
- [KDexFunctionStatus](#kdexfunctionstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `observedGeneration` _integer_ | observedGeneration is the most recent generation observed for this resource. It corresponds to the<br />resource's generation, which is updated on mutation by the API Server. |  | Optional: \{\} <br /> |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#condition-v1-meta) array_ | conditions represent the current state of the resource.<br />Each condition has a unique type and reflects the status of a specific aspect of the resource.<br />Standard condition types include:<br />- "Progressing": the resource is being created or updated<br />- "Ready": the resource is fully functional<br />- "Degraded": the resource failed to reach or maintain its desired state<br />The status of each condition is one of True, False, or Unknown. |  | Optional: \{\} <br /> |
| `attributes` _object (keys:string, values:string)_ | attributes hold state of the resource as key/value pairs. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
- [KDexClusterPageArchetype](#kdexclusterpagearchetype)
- [KDexPageArchetype](#kdexpagearchetype)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the structure of an HTML page. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `defaultFooterRef` _[KDexObjectReference](#kdexobjectreference)_ | defaultFooterRef is an optional reference to a KDexPageFooter resource. If not specified, no footer will be displayed. Use the `.Footer` property to position its content in the template. |  | Optional: \{\} <br /> |
| `defaultHeaderRef` _[KDexObjectReference](#kdexobjectreference)_ | defaultHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, no header will be displayed. Use the `.Header` property to position its content in the template. |  | Optional: \{\} <br /> |
| `defaultNavigationRefs` _object (keys:string, values:[KDexObjectReference](#kdexobjectreference))_ | defaultNavigationRefs is an optional map of keyed navigation object references. Use `.Navigation.<key>` to position the navigation's content in the template. When not empty, the 'main' key must be specified. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
| `contentEntries` _[ContentEntry](#contententry) array_ | contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or KDexApp references. |  | MaxItems: 32 <br />MinItems: 1 <br />Required: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this binding is for. |  | Required: \{\} <br /> |
| `label` _string_ | label is the value used in menus and page titles before localization occurs (or when no translation exists for the current language). |  | MaxLength: 256 <br />MinLength: 3 <br />Required: \{\} <br /> |
| `tags` _string array_ | Tags are used for grouping and searching functions. |  | Optional: \{\} <br /> |
| `contact` _[ContactInfo](#contactinfo)_ | Contact provides contact information for the function's owner. |  | Optional: \{\} <br /> |
| `navigationHints` _[NavigationHints](#navigationhints)_ | navigationHints are optional navigation properties that if omitted result in the page being hidden from the navigation. |  | Optional: \{\} <br /> |
| `overrideFooterRef` _[KDexObjectReference](#kdexobjectreference)_ | overrideFooterRef is an optional reference to a KDexPageFooter resource. If not specified, the footer from the archetype will be used. |  | Optional: \{\} <br /> |
| `overrideHeaderRef` _[KDexObjectReference](#kdexobjectreference)_ | overrideHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, the header from the archetype will be used. |  | Optional: \{\} <br /> |
| `overrideNavigationRefs` _object (keys:string, values:[KDexObjectReference](#kdexobjectreference))_ | overrideNavigationRefs is an optional map of keyed navigation object references. When not empty, the 'main' key must be specified. These navigations will be merged with the navigations from the archetype. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |
| `pageArchetypeRef` _[KDexObjectReference](#kdexobjectreference)_ | pageArchetypeRef is a reference to the KDexPageArchetype that this binding is for. |  | Required: \{\} <br /> |
| `parentPageRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | parentPageRef is a reference to the KDexPageBinding bellow which this page will appear in the main navigation. If not set, the page will be placed in the top level of the navigation. |  | Optional: \{\} <br /> |
| `basePath` _string_ | basePath is the shortest path by which the page may be accessed. It must not contain path parameters. This path will be used in site navigation. This path is subject to being prefixed for localization by `/\{l10n\}` and will be when the user selects a non-default language. |  | Pattern: `^/` <br />Required: \{\} <br /> |
| `patternPath` _string_ | patternPath, which must be prefixed by BasePath, is an extension of basePath that adds pattern matching as defined by https://pkg.go.dev/net/http#hdr-Patterns-ServeMux. This path is subject to being prefixed for localization by `/\{l10n\}` such as when the user selects a non-default language. |  | Optional: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |
| `security` _[SecurityRequirement](#securityrequirement)_ | Optional security requirements that override top-level security. |  |  |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
- [KDexClusterPageFooter](#kdexclusterpagefooter)
- [KDexPageFooter](#kdexpagefooter)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the content of an App Server page footer section. Use the `.Footer` property to position its content in the template. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
- [KDexClusterPageHeader](#kdexclusterpageheader)
- [KDexPageHeader](#kdexpageheader)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the content of an App Server page header section. Use the `.Header` property to position its content in the template. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
- [KDexClusterPageNavigation](#kdexclusterpagenavigation)
- [KDexPageNavigation](#kdexpagenavigation)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _string_ | content is a go string template that defines the content of an App Server page navigation. Use the `.Navigation["<name>"]` property to position its content in the template. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |


#### KDexScope



KDexScope is the Schema for the kdexscopes API



_Appears in:_
- [KDexScopeList](#kdexscopelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexScope` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexScopeSpec](#kdexscopespec)_ | spec defines the desired state of KDexScope |  |  |


#### KDexScopeBinding



KDexScopeBinding is the Schema for the kdexscopebindings API



_Appears in:_
- [KDexScopeBindingList](#kdexscopebindinglist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexScopeBinding` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KDexScopeBindingSpec](#kdexscopebindingspec)_ | spec defines the desired state of KDexScopeBinding |  |  |


#### KDexScopeBindingList



KDexScopeBindingList contains a list of KDexScopeBinding





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexScopeBindingList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexScopeBinding](#kdexscopebinding) array_ |  |  |  |


#### KDexScopeBindingSpec



KDexScopeBindingSpec defines the desired state of KDexScopeBinding



_Appears in:_
- [KDexScopeBinding](#kdexscopebinding)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `email` _string_ | email is the email address of the subject, used for local fallback lookup or metadata. |  | Optional: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this binding is for. |  | Required: \{\} <br /> |
| `scopes` _string array_ | scopes is a list of internal scopes bound to this subject. |  | MinItems: 1 <br />Required: \{\} <br /> |
| `secretRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | secretRef is an optional reference to a secret that contains keys that map to subject and<br />the value is the password. As such the secret can be mapped to multiple KDexScopeBinding.<br />This simple fallback is not intended for large scale production use. Thought it may be used for administration. |  | Optional: \{\} <br /> |
| `subject` _string_ | subject is the subject identifier. It should be from the OIDC provider (e.g. Google).<br />However, if the secretRef is set then it contains a local identity managed<br />through the Secret. |  | MinLength: 5 <br />Required: \{\} <br /> |


#### KDexScopeList



KDexScopeList contains a list of KDexScope





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexScopeList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexScope](#kdexscope) array_ |  |  |  |


#### KDexScopeSpec



KDexScopeSpec defines the desired state of KDexScope



_Appears in:_
- [KDexScope](#kdexscope)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this binding is for. |  | Required: \{\} <br /> |
| `rules` _[PolicyRule](#policyrule) array_ | Rules holds all the PolicyRules for this KDexScope |  | MinItems: 1 <br />Required: \{\} <br /> |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
- [KDexClusterScriptLibrary](#kdexclusterscriptlibrary)
- [KDexScriptLibrary](#kdexscriptlibrary)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `packageReference` _[PackageReference](#packagereference)_ | packageReference specifies the name and version of an NPM package that contains the script. The package.json must describe an ES module. |  | Optional: \{\} <br /> |
| `scripts` _[ScriptDef](#scriptdef) array_ | scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents. |  | MaxItems: 8 <br />Optional: \{\} <br /> |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | imagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.<br />More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod |  | Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/_' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/_.+` <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |
| `serverImage` _string_ | serverImage is the name of Backend image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `serverImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the Backend server image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `staticImage` _string_ | staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `staticImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
- [KDexClusterTheme](#kdexclustertheme)
- [KDexTheme](#kdextheme)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `assets` _[Assets](#assets)_ | assets is a set of elements that define a portable set of design rules. |  | MaxItems: 32 <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | imagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.<br />More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod |  | Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/_' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/_.+` <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |
| `serverImage` _string_ | serverImage is the name of Backend image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `serverImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the Backend server image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `staticImage` _string_ | staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `staticImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
- [KDexClusterTranslation](#kdexclustertranslation)
- [KDexInternalTranslationSpec](#kdexinternaltranslationspec)
- [KDexTranslation](#kdextranslation)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `translations` _[Translation](#translation) array_ | translations is an array of objects where each one specifies a language (lang) and a map (keysAndValues) consisting of key/value pairs. If the lang property is not unique in the array and its keysAndValues map contains the same keys, the last one takes precedence. |  | MaxItems: 32 <br />MinItems: 1 <br />Required: \{\} <br /> |


#### KDexUtilityPage



KDexUtilityPage is the Schema for the kdexutilitypages API

A KDexUtilityPage defines a utility page (Announcement, Error, Login) that can be referenced by a KDexHost.
It shares much of its structure with KDexPageBinding but is specialized for system-level pages that do not
necessarily sit within the standard site navigation tree.



_Appears in:_
- [KDexUtilityPageList](#kdexutilitypagelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexUtilityPage` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexUtilityPageSpec](#kdexutilitypagespec)_ | spec defines the desired state of KDexUtilityPage |  | Required: \{\} <br /> |


#### KDexUtilityPageList



KDexUtilityPageList contains a list of KDexUtilityPage





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexUtilityPageList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexUtilityPage](#kdexutilitypage) array_ |  |  |  |


#### KDexUtilityPageSpec



KDexUtilityPageSpec defines the desired state of KDexUtilityPage



_Appears in:_
- [KDexClusterUtilityPage](#kdexclusterutilitypage)
- [KDexInternalUtilityPageSpec](#kdexinternalutilitypagespec)
- [KDexUtilityPage](#kdexutilitypage)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[KDexUtilityPageType](#kdexutilitypagetype)_ | type indicates the purpose of this utility page. |  | Enum: [Announcement Error Login] <br />Required: \{\} <br /> |
| `contentEntries` _[ContentEntry](#contententry) array_ | contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or KDexApp references. |  | MaxItems: 8 <br />MinItems: 1 <br />Required: \{\} <br /> |
| `overrideFooterRef` _[KDexObjectReference](#kdexobjectreference)_ | overrideFooterRef is an optional reference to a KDexPageFooter resource. If not specified, the footer from the archetype will be used. |  | Optional: \{\} <br /> |
| `overrideHeaderRef` _[KDexObjectReference](#kdexobjectreference)_ | overrideHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, the header from the archetype will be used. |  | Optional: \{\} <br /> |
| `overrideNavigationRefs` _object (keys:string, values:[KDexObjectReference](#kdexobjectreference))_ | overrideNavigationRefs is an optional map of keyed navigation object references. When not empty, the 'main' key must be specified. These navigations will be merged with the navigations from the archetype. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |
| `pageArchetypeRef` _[KDexObjectReference](#kdexobjectreference)_ | pageArchetypeRef is a reference to the KDexPageArchetype that this binding is for. |  | Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |


#### KDexUtilityPageType

_Underlying type:_ _string_

KDexUtilityPageType defines the type of utility page.

_Validation:_
- Enum: [Announcement Error Login]

_Appears in:_
- [KDexInternalUtilityPageSpec](#kdexinternalutilitypagespec)
- [KDexUtilityPageSpec](#kdexutilitypagespec)

| Field | Description |
| --- | --- |
| `Announcement` | AnnouncementUtilityPageType represents an announcement page.<br /> |
| `Error` | ErrorUtilityPageType represents an error page.<br /> |
| `Login` | LoginUtilityPageType represents a login page.<br /> |


#### LocalSecretWithKeyReference







_Appears in:_
- [JWT](#jwt)
- [OIDCProvider](#oidcprovider)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `keyProperty` _string_ | keyProperty is the property from which to extract a value from the secret |  | Required: \{\} <br /> |
| `secretRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | secretRef is a reference to a secret in the same namespace as the referrer. |  | Required: \{\} <br /> |


#### MappingRule







_Appears in:_
- [OIDCProvider](#oidcprovider)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `required` _boolean_ | required indicates that if the rule fails to produce a value token generation should fail as well | false | Optional: \{\} <br /> |
| `expession` _string_ | expession is CEL program to compute a transformation of claims from the OIDC token. |  | Required: \{\} <br /> |
| `target` _string_ | target is a nested property path to where the result will be attached to the claims structure |  | Required: \{\} <br /> |


#### Metadata



KDexFunctionMetadata defines the metadata for the function.



_Appears in:_
- [KDexFunctionMetadata](#kdexfunctionmetadata)
- [KDexPageBindingSpec](#kdexpagebindingspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `tags` _string array_ | Tags are used for grouping and searching functions. |  | Optional: \{\} <br /> |
| `contact` _[ContactInfo](#contactinfo)_ | Contact provides contact information for the function's owner. |  | Optional: \{\} <br /> |


#### ModulePolicy

_Underlying type:_ _string_

ModulePolicy defines the policy for the use of JavaScript Modules.

_Validation:_
- Enum: [ExternalDependencies Loose ModulesRequired Strict]

_Appears in:_
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)

| Field | Description |
| --- | --- |
| `Loose` | LooseModulePolicy means that a) JavaScript references are not required to be JavaScript modules and b) JavaScript references may contain embed dependencies.<br /> |
| `ExternalDependencies` | ExternalDependenciesModulePolicy means that a) JavaScript references are not required to be JavaScript modules and b) JavaScript references may not contain embed dependencies.<br /> |
| `ModulesRequired` | ModulesRequiredModulePolicy means that a) JavaScript references are required to be JavaScript modules and b) JavaScript references may contain embed dependencies.<br /> |
| `Strict` | StrictModulePolicy means that a) JavaScript references are required to be JavaScript modules and b) JavaScript references may not contain embed dependencies.<br /> |


#### NavigationHints







_Appears in:_
- [KDexPageBindingSpec](#kdexpagebindingspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `icon` _string_ | icon is the name of the icon to display next to the menu entry for this page. |  | Optional: \{\} <br /> |
| `weight` _[Quantity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#quantity-resource-api)_ | weight is a property that influences the position of the page menu entry. Items at each level are sorted first by ascending weight and then ascending lexicographically. |  | Optional: \{\} <br /> |


#### OIDCProvider







_Appears in:_
- [Auth](#auth)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `clientID` _string_ | clientID is the id assigned by the provider to this application. |  | Required: \{\} <br /> |
| `clientSecretRef` _[LocalSecretWithKeyReference](#localsecretwithkeyreference)_ | clientSecretRef is a reference to a secret in the host's namespace that holds the client_secret assigned to this application by the OIDC provider. |  | Required: \{\} <br /> |
| `mappers` _[MappingRule](#mappingrule) array_ | mappers is an array of CEL expressions for extracting custom claims from the OIDC token and adding the results to the local token. |  | MaxItems: 16 <br />Optional: \{\} <br /> |
| `oidcProviderURL` _string_ | oidcProviderURL is the well known URL of the OIDC provider. |  | Required: \{\} <br /> |
| `scopes` _string array_ | scopes is an array of additional scopes that will be requested from the provider. |  | Optional: \{\} <br /> |


#### PackageReference



PackageReference specifies the name and version of an NPM package. Prefereably the package should be available from
the public npm registry. If the package is not available from the public npm registry, a secretRef should be provided
to authenticate to the npm registry. That package must contain an ES module for use in the browser.



_Appears in:_
- [KDexAppSpec](#kdexappspec)
- [KDexInternalPackageReferencesSpec](#kdexinternalpackagereferencesspec)
- [KDexScriptLibrarySpec](#kdexscriptlibraryspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | name contains a scoped npm package name. |  | Required: \{\} <br /> |
| `version` _string_ | version contains a specific npm package version. |  | Required: \{\} <br /> |
| `exportMapping` _string_ | exportMapping is a mapping of the module's exports that will be used when the module import is written. e.g. `import [exportMapping] from [module_name];`. If exportMapping is not provided the module will be written as `import [module_name];` |  | Optional: \{\} <br /> |
| `secretRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | secretRef is a reference to a secret containing authentication credentials for the NPM registry that holds the package. |  | Optional: \{\} <br /> |


#### PathItem







_Appears in:_
- [API](#api)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `connect` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg)_ |  |  | Optional: \{\} <br />Type: object <br /> |
| `delete` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg)_ |  |  | Optional: \{\} <br />Type: object <br /> |
| `description` _string_ |  |  | Optional: \{\} <br /> |
| `get` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg)_ |  |  | Optional: \{\} <br />Type: object <br /> |
| `head` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg)_ |  |  | Optional: \{\} <br />Type: object <br /> |
| `options` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg)_ |  |  | Optional: \{\} <br />Type: object <br /> |
| `parameters` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg) array_ |  |  | Optional: \{\} <br /> |
| `patch` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg)_ |  |  | Optional: \{\} <br />Type: object <br /> |
| `post` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg)_ |  |  | Optional: \{\} <br />Type: object <br /> |
| `put` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg)_ |  |  | Optional: \{\} <br />Type: object <br /> |
| `summary` _string_ |  |  | Optional: \{\} <br /> |
| `trace` _[RawExtension](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#rawextension-runtime-pkg)_ |  |  | Optional: \{\} <br />Type: object <br /> |


#### Paths







_Appears in:_
- [KDexPageBindingSpec](#kdexpagebindingspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `basePath` _string_ | basePath is the shortest path by which the page may be accessed. It must not contain path parameters. This path will be used in site navigation. This path is subject to being prefixed for localization by `/\{l10n\}` and will be when the user selects a non-default language. |  | Pattern: `^/` <br />Required: \{\} <br /> |
| `patternPath` _string_ | patternPath, which must be prefixed by BasePath, is an extension of basePath that adds pattern matching as defined by https://pkg.go.dev/net/http#hdr-Patterns-ServeMux. This path is subject to being prefixed for localization by `/\{l10n\}` such as when the user selects a non-default language. |  | Optional: \{\} <br /> |


#### PolicyRule



PolicyRule holds information that describes a policy rule, but does not
contain information about who the rule applies to or which namespace the
rule applies to.



_Appears in:_
- [KDexScopeSpec](#kdexscopespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `resourceNames` _string array_ | resourceNames is an optional allow list of names that the rule applies to. An empty set means that everything is allowed. |  |  |
| `resources` _string array_ | resources is a list of resources this rule applies to. '*' represents all resources. |  | MinItems: 1 <br />Required: \{\} <br /> |
| `verbs` _string array_ | verbs is a list of verbs that apply to ALL the resources contained in this rule. '*' represents all verbs. |  | MinItems: 1 <br />Required: \{\} <br /> |


#### Routing



Routing defines the desired routing configuration for the host.



_Appears in:_
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `domains` _string array_ | domains are the names by which this host is addressed. The first domain listed is the preferred domain. The domains may contain wildcard prefix in the form '*.'. Longest match always wins. |  | MinItems: 1 <br />Required: \{\} <br /> |
| `ingressClassName` _string_ | ingressClassName is the name of an IngressClass cluster resource. Ingress<br />controller implementations use this field to know whether they should be<br />serving this Ingress resource, by a transitive connection<br />(controller -> IngressClass -> Ingress resource). Although the<br />`kubernetes.io/ingress.class` annotation (simple constant name) was never<br />formally defined, it was widely supported by Ingress controllers to create<br />a direct binding between Ingress controller and Ingress resources. Newly<br />created Ingress resources should prefer using the field. However, even<br />though the annotation is officially deprecated, for backwards compatibility<br />reasons, ingress controllers should still honor that annotation if present. |  | Optional: \{\} <br /> |
| `strategy` _[RoutingStrategy](#routingstrategy)_ | strategy is the routing strategy to use. If not specified Ingress is assumed. | Ingress | Enum: [Ingress HTTPRoute] <br />Optional: \{\} <br /> |
| `tls` _[TLSSpec](#tlsspec)_ | tls is the TLS configuration for the host. |  | Optional: \{\} <br /> |


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


#### ScalingConfig



ScalingConfig defines scaling parameters.



_Appears in:_
- [KDexFunctionExec](#kdexfunctionexec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `maxReplicas` _integer_ | MaxReplicas is the maximum number of replicas. |  | Minimum: 1 <br />Optional: \{\} <br /> |
| `minReplicas` _integer_ | MinReplicas is the minimum number of replicas. |  | Minimum: 0 <br />Optional: \{\} <br /> |


#### ScriptDef







_Appears in:_
- [KDexAppSpec](#kdexappspec)
- [KDexScriptLibrarySpec](#kdexscriptlibraryspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `script` _string_ | script is the content that will be added to a `<script>` element when rendered. |  | Optional: \{\} <br /> |
| `scriptSrc` _string_ | scriptSrc is a value for a `<script>` `src` attribute. It must be either and absolute URL with a protocol and host<br />or it must be relative to the `ingressPath` field of the specified Backend. |  | Optional: \{\} <br /> |
| `footScript` _boolean_ | footScript is a flag for script or scriptSrc that indicates if the tag should be added in the head of the page or at the foot. The default is false (add to head). To add the script to the foot of the page set footScript to true. | false | Optional: \{\} <br /> |
| `attributes` _object (keys:string, values:string)_ | attributes are key/value pairs that will be added to the element when rendered. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |


#### SecurityRequirement

_Underlying type:_ _object_





_Appears in:_
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)
- [KDexPageBindingSpec](#kdexpagebindingspec)



#### StubDetails



StubDetails contains stub information.



_Appears in:_
- [KDexFunctionStatus](#kdexfunctionstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `filePath` _string_ | FilePath is the path to the generated function file. |  | Optional: \{\} <br /> |
| `language` _string_ | Language is the programming language of the stub. |  | Optional: \{\} <br /> |
| `sourceImage` _string_ | SourceImage is the OCI artifact reference where the stub code was pushed. |  | Optional: \{\} <br /> |




#### TLSSpec



TLSSpec defines the desired state of TLS for a host.



_Appears in:_
- [Routing](#routing)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `secretRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | secretRef is a reference to a secret containing a TLS certificate and key for the domains specified on the host. |  | Required: \{\} <br /> |


#### Translation







_Appears in:_
- [KDexInternalTranslationSpec](#kdexinternaltranslationspec)
- [KDexTranslationSpec](#kdextranslationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `lang` _string_ | lang is a string containing a BCP 47 language tag that identifies the set of translations.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag. |  | Required: \{\} <br /> |
| `keysAndValues` _object (keys:string, values:string)_ | keysAndValues is a map of key=/value pairs where the key is the identifier and the value is the translation of that key in the language specified by the lang property. |  | MaxProperties: 256 <br />MinProperties: 1 <br />Required: \{\} <br /> |


#### UtilityPages



UtilityPages defines the utility pages for a host.



_Appears in:_
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `announcementRef` _[KDexObjectReference](#kdexobjectreference)_ | announcementRef is a reference to a KDexUtilityPage or KDexClusterUtilityPage that defines the announcement page. |  | Optional: \{\} <br /> |
| `errorRef` _[KDexObjectReference](#kdexobjectreference)_ | errorRef is a reference to a KDexUtilityPage or KDexClusterUtilityPage that defines the error page. |  | Optional: \{\} <br /> |
| `loginRef` _[KDexObjectReference](#kdexobjectreference)_ | loginRef is a reference to a KDexUtilityPage or KDexClusterUtilityPage that defines the login page. |  | Optional: \{\} <br /> |


