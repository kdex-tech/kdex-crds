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
- [KDexClusterFaaSAdaptor](#kdexclusterfaasadaptor)
- [KDexClusterFaaSAdaptorList](#kdexclusterfaasadaptorlist)
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
- [KDexFaaSAdaptor](#kdexfaasadaptor)
- [KDexFaaSAdaptorList](#kdexfaasadaptorlist)
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
- [KDexRole](#kdexrole)
- [KDexRoleBinding](#kdexrolebinding)
- [KDexRoleBindingList](#kdexrolebindinglist)
- [KDexRoleList](#kdexrolelist)
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





_Validation:_
- ExactlyOneOf: [linkHref metaId style]

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
- ExactlyOneOf: [linkHref metaId style]
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
| `anonymousEntitlements` _string array_ | anonymousEntitlements is an array of entitlements granted in anonymous (not logged in) access scenarios.<br />In the spirit of least privilege security no entitlements are granted by default. However, in order to make<br />a host's pages generally accessible the scope `page:read` should be granted. |  | Optional: \{\} <br /> |
| `claimMappings` _MappingRule array_ | claimMappings is an array of CEL expressions for extracting custom claims from<br />identity sources and mapping the results onto the Primary Access Token (PAT).<br />This is used to map OIDC claims but can also be used with external data<br />sources like LDAP or others via identity integration. |  | MaxItems: 16 <br />Optional: \{\} <br /> |
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
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/-/' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/-/.+` <br /> |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is an optional list of environment variables to set in the container. |  | Optional: \{\} <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |
| `serverImage` _string_ | serverImage is the name of Backend image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `serverImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the Backend server image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `staticImage` _string_ | staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `staticImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |


#### Builder







_Appears in:_
- [KDexFaaSAdaptorSpec](#kdexfaasadaptorspec)
- [Source](#source)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `builderRef` _[KDexObjectReference](#kdexobjectreference)_ | builderRef is a reference to the kpack.io/v1alpha2/Builder or kpack.io/v1alpha2/ClusterBuilder to use for building the image. |  | Required: \{\} <br /> |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is the environment variables to set in the builder. |  | Optional: \{\} <br /> |
| `languages` _string array_ | Languages is a list of languages that this builder supports. |  | MinItems: 1 <br />Required: \{\} <br /> |
| `name` _string_ | Name is the builder name (e.g., tiny, base, full). |  | Required: \{\} <br /> |
| `serviceAccountName` _string_ | serviceAccountName is the name of the service account to use for building the image. |  | Optional: \{\} <br /> |








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





_Validation:_
- ExactlyOneOf: [appRef rawHTML]

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


#### Deployer







_Appears in:_
- [KDexFaaSAdaptorSpec](#kdexfaasadaptorspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `args` _string array_ | args is an optional array of arguments that will be passed to the generator command. |  | Optional: \{\} <br /> |
| `command` _string array_ | command is an optional array that contains the code generator command and any flags necessary. |  | Optional: \{\} <br /> |
| `image` _string_ | image is the image to use for deploying executables into a FaaS runtime. |  | Required: \{\} <br /> |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is the environment variables to set in the deployer. |  | Optional: \{\} <br /> |
| `serviceAccountName` _string_ | serviceAccountName is the name of the service account to use for deploying executables into a FaaS runtime. |  | Optional: \{\} <br /> |


#### Executable







_Appears in:_
- [FunctionOrigin](#functionorigin)
- [KDexFunctionStatus](#kdexfunctionstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | image is a reference to executable artifact. In most cases this will be a Docker image. In some other cases<br />it may be an artifact native to FaaS Adaptor's target runtime. |  | Optional: \{\} <br /> |
| `scaling` _[ScalingConfig](#scalingconfig)_ | Scaling allows configuration for min/max replicas and autoscaler type. |  | Optional: \{\} <br /> |


#### FunctionOrigin



FunctionOrigin defines the origin of the function implementation.
There are four possible ways to obtain a deployable function:
1. Executable: A pre-built container image or VM image.
2. Source: A reference to source code that will be compiled and built into a container image or VM image.
3. Generator: A configuration for a code generator that will produce source code stubs for the selected language.
4. Nothing: A code generator config will be derived from the defaults provided by the FaaS Adaptor.

_Validation:_
- AtMostOneOf: [executable generator source]

_Appears in:_
- [KDexFunctionSpec](#kdexfunctionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `executable` _[Executable](#executable)_ | executable is a reference to a pre-built container image or VM image. |  | Optional: \{\} <br /> |
| `generator` _[Generator](#generator)_ | generator holds the values to configure and execute the code generator. |  | Optional: \{\} <br /> |
| `source` _[Source](#source)_ | source contains source code location information. |  | Optional: \{\} <br /> |


#### Generator







_Appears in:_
- [FunctionOrigin](#functionorigin)
- [KDexFaaSAdaptorSpec](#kdexfaasadaptorspec)
- [KDexFunctionStatus](#kdexfunctionstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `args` _string array_ | args is an optional array of arguments that will be passed to the generator command. |  | Optional: \{\} <br /> |
| `command` _string array_ | command is an optional array that contains the code generator command and any flags necessary. |  | Optional: \{\} <br /> |
| `entrypoint` _string_ | Entrypoint is the specific function handler/method to execute. |  | Optional: \{\} <br /> |
| `git` _[Git](#git)_ | git is the configuration for the Git repository where generated code will be committed to a branch. |  | Required: \{\} <br /> |
| `image` _string_ | image is the image containing the generator implementation; cli or scripts. |  | Required: \{\} <br /> |
| `language` _string_ | Language is the programming language of the function (e.g., go, python, nodejs). |  | Required: \{\} <br /> |
| `serviceAccountName` _string_ | serviceAccountName is the name of the service account to use for the generator job. |  | Optional: \{\} <br /> |


#### Git







_Appears in:_
- [Generator](#generator)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `functionSubDirectory` _string_ | functionSubDirectory is the optional path to a subdirectory in the repository in which generated code will be placed. | . | Optional: \{\} <br /> |
| `image` _string_ | image is the name of the container image to run for git. |  | Required: \{\} <br /> |
| `committerEmail` _string_ | committerEmail is the email address that will be used for git commits. |  | Required: \{\} <br /> |
| `committerName` _string_ | committerName is the name that will be used for git commits. |  | Required: \{\} <br /> |


#### JWT







_Appears in:_
- [Auth](#auth)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `cookieName` _string_ | cookieName is the name of the Cookie in which the JWT token will be stored. (default is "auth_token") | auth_token | Optional: \{\} <br /> |
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
| `scripts` _[ScriptDef](#scriptdef) array_ | scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents. |  | ExactlyOneOf: [script scriptSrc] <br />MaxItems: 8 <br />Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/-/' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/-/.+` <br /> |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is an optional list of environment variables to set in the container. |  | Optional: \{\} <br /> |
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
| `spec` _[KDexAppSpec](#kdexappspec)_ | spec defines the desired state of KDexClusterApp |  | Required: \{\} <br /> |


#### KDexClusterAppList



KDexClusterAppList contains a list of KDexClusterApp





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterAppList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterApp](#kdexclusterapp) array_ |  |  |  |


#### KDexClusterFaaSAdaptor



KDexClusterFaaSAdaptor is the Schema for the kdexclusterfaasadaptors API



_Appears in:_
- [KDexClusterFaaSAdaptorList](#kdexclusterfaasadaptorlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterFaaSAdaptor` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexFaaSAdaptorSpec](#kdexfaasadaptorspec)_ | spec defines the desired state of KDexClusterFaaSAdaptor |  | Required: \{\} <br /> |


#### KDexClusterFaaSAdaptorList



KDexClusterFaaSAdaptorList contains a list of KDexClusterFaaSAdaptor





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterFaaSAdaptorList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexClusterFaaSAdaptor](#kdexclusterfaasadaptor) array_ |  |  |  |


#### KDexClusterPageArchetype



KDexClusterPageArchetype is the Schema for the kdexclusterpagearchetypes API



_Appears in:_
- [KDexClusterPageArchetypeList](#kdexclusterpagearchetypelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexClusterPageArchetype` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexPageArchetypeSpec](#kdexpagearchetypespec)_ | spec defines the desired state of KDexClusterPageArchetype |  | Required: \{\} <br /> |


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
| `spec` _[KDexPageFooterSpec](#kdexpagefooterspec)_ | spec defines the desired state of KDexClusterPageFooter |  | Required: \{\} <br /> |


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
| `spec` _[KDexPageHeaderSpec](#kdexpageheaderspec)_ | spec defines the desired state of KDexClusterPageHeader |  | Required: \{\} <br /> |


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
| `spec` _[KDexPageNavigationSpec](#kdexpagenavigationspec)_ | spec defines the desired state of KDexClusterPageNavigation |  | Required: \{\} <br /> |


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
| `spec` _[KDexScriptLibrarySpec](#kdexscriptlibraryspec)_ | spec defines the desired state of KDexClusterScriptLibrary |  | Required: \{\} <br /> |


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
| `spec` _[KDexThemeSpec](#kdexthemespec)_ | spec defines the desired state of KDexClusterTheme |  | Required: \{\} <br /> |


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
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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


#### KDexFaaSAdaptor



KDexFaaSAdaptor is the Schema for the kdexfaasadaptors API



_Appears in:_
- [KDexFaaSAdaptorList](#kdexfaasadaptorlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexFaaSAdaptor` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexFaaSAdaptorSpec](#kdexfaasadaptorspec)_ | spec defines the desired state of KDexFaaSAdaptor |  | Required: \{\} <br /> |


#### KDexFaaSAdaptorList



KDexFaaSAdaptorList contains a list of KDexFaaSAdaptor





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexFaaSAdaptorList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexFaaSAdaptor](#kdexfaasadaptor) array_ |  |  |  |


#### KDexFaaSAdaptorSpec



KDexFaaSAdaptorSpec defines the desired state of KDexFaaSAdaptor



_Appears in:_
- [KDexClusterFaaSAdaptor](#kdexclusterfaasadaptor)
- [KDexFaaSAdaptor](#kdexfaasadaptor)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `builders` _[Builder](#builder) array_ | Builders is a list of builder configurations. |  | MinItems: 1 <br /> |
| `defaultBuilderGenerator` _string_ | DefaultBuilderGenerator is the default builder/generator combination to use for functions that do not specify a builder or generator.<br />The format is "<builder>/<generator>" (e.g., "tiny/go"). |  | Pattern: `^\w+/\w+$` <br />Required: \{\} <br /> |
| `deployer` _[Deployer](#deployer)_ | Deployer is the configuration for the deployer. |  | Required: \{\} <br /> |
| `generators` _[Generator](#generator) array_ | Generators is a list of provider-specific generator configurations. |  | MinItems: 1 <br /> |
| `observer` _[Observer](#observer)_ | Observer is the configuration for the observer. |  | Optional: \{\} <br /> |
| `provider` _string_ | Provider is the type of FaaS provider (e.g., "knative", "openfaas", "lambda"). |  | Enum: [knative openfaas lambda azure-functions google-cloud-functions] <br />Required: \{\} <br /> |




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
| `tags` _[Tag](#tag) array_ | Tags are used for grouping and searching functions. |  | MaxItems: 16 <br />Optional: \{\} <br /> |
| `contact` _[ContactInfo](#contactinfo)_ | Contact provides contact information for the function's owner. |  | Optional: \{\} <br /> |


#### KDexFunctionSpec



KDexFunctionSpec defines the desired state of KDexFunction



_Appears in:_
- [KDexFunction](#kdexfunction)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `api` _[API](#api)_ | api defines the OpenAPI contract for the function.<br />See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#path-item-object<br />See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#schema-object<br />The supported fields from 'path item object' are: summary, description, get, put, post, delete, options, head, patch, trace, parameters, and responses.<br />The field 'schemas' of type map[string]schema whose values are defined by 'schema object' is supported and can be referenced throughout operation definitions. References must be in the form "#/components/schemas/<name>". |  | Required: \{\} <br /> |
| `claimMappings` _MappingRule array_ | claimMappings is an array of CEL expressions for extracting custom claims<br />from the current authorization context onto the Function Access Token (FAT).<br />This can be used to map Function specific claims like tenant, department_id,<br />strip_customer_id, etc. to the FAT. |  | MaxItems: 16 <br />Optional: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this translation belongs to. |  | Required: \{\} <br /> |
| `metadata` _[KDexFunctionMetadata](#kdexfunctionmetadata)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `origin` _[FunctionOrigin](#functionorigin)_ | origin defines the origin of the function implementation. |  | AtMostOneOf: [executable generator source] <br />Optional: \{\} <br /> |


#### KDexFunctionState

_Underlying type:_ _string_

KDexFunctionState reflects the current state of a KDexFunction.

_Validation:_
- Enum: [Pending OpenAPIValid BuildValid SourceAvailable ExecutableAvailable FunctionDeployed Ready]

_Appears in:_
- [KDexFunctionStatus](#kdexfunctionstatus)

| Field | Description |
| --- | --- |
| `Pending` | 1. KDexFunctionStatePending indicates the function is pending action.<br /> |
| `OpenAPIValid` | 2. KDexFunctionStateOpenAPIValid indicates the OpenAPI spec is valid.<br /> |
| `BuildValid` | 3. KDexFunctionStateBuildValid indicates the build configuration is valid.<br /> |
| `SourceAvailable` | 4. KDexFunctionStateSourceAvailable indicates the source code is available.<br /> |
| `ExecutableAvailable` | 5. KDexFunctionStateExecutableAvailable indicates the executable container is available for provisioning.<br /> |
| `FunctionDeployed` | 6. KDexFunctionStateFunctionDeployed indicates the function has been deployed to the FaaS runtime.<br /> |
| `Ready` | 7. KDexFunctionStateReady indicates the function is ready for invocation.<br /> |




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
| `assets` _[Assets](#assets)_ | assets is a set of elements that define a host specific HTML instructions (e.g. favicon, site logo, charset). |  | ExactlyOneOf: [linkHref metaId style] <br />MaxItems: 32 <br /> |
| `auth` _[Auth](#auth)_ | auth holds the host's authentication configuration. |  | Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/-/' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/-/.+` <br /> |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is an optional list of environment variables to set in the container. |  | Optional: \{\} <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |
| `serverImage` _string_ | serverImage is the name of Backend image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `serverImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the Backend server image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `staticImage` _string_ | staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `staticImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `brandName` _string_ | brandName is the name used when rendering pages belonging to the host. For example, it may be used as alt text for the logo displayed in the page header. |  | Required: \{\} <br /> |
| `defaultLang` _string_ | defaultLang is a string containing a BCP 47 language tag.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.<br />When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'. |  | Optional: \{\} <br /> |
| `devMode` _boolean_ | devMode is a boolean that enables development features like the Request Sniffer. |  | Optional: \{\} <br /> |
| `faasAdaptorRef` _[KDexObjectReference](#kdexobjectreference)_ | faasAdaptorRef is an optional reference to the FaaS Adaptor that will drive KDexFunction generation of code and deployment for this host. If not specified the default will be used. |  | Optional: \{\} <br /> |
| `faviconSVGTemplate` _string_ | faviconSVGTemplate contains SVG code marked up with go string template to which will be passed the render.TemplateData holding other host details. The rendered output will be cached and served at "/favicon.ico" as "image/svg+xml". |  | Optional: \{\} <br /> |
| `modulePolicy` _[ModulePolicy](#modulepolicy)_ | modulePolicy defines the policy for JavaScript references in KDexApp, KDexTheme and KDexScriptLibrary resources. When not specified the policy is Strict<br />A Host must not accept JavaScript references which do not comply with the specified policy. | Strict | Enum: [ExternalDependencies Loose ModulesRequired Strict] <br />Optional: \{\} <br /> |
| `openapi` _[OpenAPI](#openapi)_ | openapi holds the configuration for the host's OpenAPI support. |  | Optional: \{\} <br /> |
| `organization` _string_ | organization is the name of the Organization to which the host belongs. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `routing` _[Routing](#routing)_ | routing defines the desired routing configuration for the host. |  | Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |
| `security` _[SecurityRequirement](#securityrequirement)_ | Optional top level security requirements. |  |  |
| `themeRef` _[KDexObjectReference](#kdexobjectreference)_ | themeRef is a reference to the theme that should apply to all pages bound to this host. |  | Optional: \{\} <br /> |
| `serviceAccountRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | serviceAccountRef is a reference to the service account used by the host to access secrets.<br />Each Secret must match one of the following:<br />- is annotated with 'kdex.dev/secret-type = auth-client' (multiple)<br />    - must contain key 'client-id' OR 'client_id'<br />    - may contain key 'public' (true\|false, default: false)<br />    - if not public, must contain key 'client-secret' OR 'client_secret'<br />    - must contain key 'redirect-uris' OR 'redirect_uris' (comma separated list)<br />    - may contain key 'allowed-grant-types' OR 'allowed_grant_types' (comma separated list)<br />    - may contain key 'allowed-scopes' OR 'allowed_scopes' (comma separated list)<br />    - may contain key 'require-pkce' OR 'require_pkce' (true\|false, default: false)<br />    - may contain key 'name'<br />    - may contain key 'description'<br />- is annotated with 'kdex.dev/secret-type = git' (single)<br />    - must contain key 'host'<br />    - must contain key 'org'<br />    - must contain key 'password'<br />    - must contain key 'repo'<br />    - must contain key 'username'<br />- is annotated with 'kdex.dev/secret-type = jwt-keys' (multiple)<br />    - must contain key 'private-key'<br />    - may be annotated with 'kdex.dev/active-key = true'<br />- is annotated with 'kdex.dev/secret-type = ldap' (single)<br />    - must contain key 'active-directory' (true\|false)<br />    - must contain key 'addr'<br />    - must contain key 'base-dn'<br />    - must contain key 'bind-dn'<br />    - must contain key 'bind-user'<br />    - must contain key 'bind-pass'<br />    - must contain key 'user-filter'<br />    - may contain key 'attributes' (comma separated list of attributes to retrieve)<br />- is annotated with 'kdex.dev/secret-type = npm' (single)<br />    - must contain key '.npmrc' (formatted as a complete .npmrc file)<br />- is annotated with 'kdex.dev/secret-type = oidc-client' (single)<br />    - must contain key 'client-id' OR 'client_id'<br />    - must contain key 'client-secret' OR 'client_secret'<br />    - may contain key 'block-key' OR 'block_key'<br />- is annotated with 'kdex.dev/secret-type = subject' (multiple)<br />    - must contain key 'sub'<br />    - must contain key 'password'<br />    - may contain arbitrary key(string)/value(string\|yaml) pairs which can be mapped to the claims using the spec.auth.claimMappings<br />- is of type 'kubernetes.io/dockerconfigjson' (multiple)<br />- is of type 'kubernetes.io/tls' (single) |  | Required: \{\} <br /> |
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
| `assets` _[Assets](#assets)_ | assets is a set of elements that define a host specific HTML instructions (e.g. favicon, site logo, charset). |  | ExactlyOneOf: [linkHref metaId style] <br />MaxItems: 32 <br /> |
| `auth` _[Auth](#auth)_ | auth holds the host's authentication configuration. |  | Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/-/' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/-/.+` <br /> |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is an optional list of environment variables to set in the container. |  | Optional: \{\} <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |
| `serverImage` _string_ | serverImage is the name of Backend image.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `serverImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the Backend server image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `staticImage` _string_ | staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.<br />More info: https://kubernetes.io/docs/concepts/containers/images |  | MinLength: 5 <br />Optional: \{\} <br /> |
| `staticImagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Policy for pulling the OCI theme image. Possible values are:<br />Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.<br />Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.<br />IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. |  | Optional: \{\} <br /> |
| `brandName` _string_ | brandName is the name used when rendering pages belonging to the host. For example, it may be used as alt text for the logo displayed in the page header. |  | Required: \{\} <br /> |
| `defaultLang` _string_ | defaultLang is a string containing a BCP 47 language tag.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.<br />When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'. |  | Optional: \{\} <br /> |
| `devMode` _boolean_ | devMode is a boolean that enables development features like the Request Sniffer. |  | Optional: \{\} <br /> |
| `faasAdaptorRef` _[KDexObjectReference](#kdexobjectreference)_ | faasAdaptorRef is an optional reference to the FaaS Adaptor that will drive KDexFunction generation of code and deployment for this host. If not specified the default will be used. |  | Optional: \{\} <br /> |
| `faviconSVGTemplate` _string_ | faviconSVGTemplate contains SVG code marked up with go string template to which will be passed the render.TemplateData holding other host details. The rendered output will be cached and served at "/favicon.ico" as "image/svg+xml". |  | Optional: \{\} <br /> |
| `modulePolicy` _[ModulePolicy](#modulepolicy)_ | modulePolicy defines the policy for JavaScript references in KDexApp, KDexTheme and KDexScriptLibrary resources. When not specified the policy is Strict<br />A Host must not accept JavaScript references which do not comply with the specified policy. | Strict | Enum: [ExternalDependencies Loose ModulesRequired Strict] <br />Optional: \{\} <br /> |
| `openapi` _[OpenAPI](#openapi)_ | openapi holds the configuration for the host's OpenAPI support. |  | Optional: \{\} <br /> |
| `organization` _string_ | organization is the name of the Organization to which the host belongs. |  | MinLength: 5 <br />Required: \{\} <br /> |
| `routing` _[Routing](#routing)_ | routing defines the desired routing configuration for the host. |  | Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |
| `security` _[SecurityRequirement](#securityrequirement)_ | Optional top level security requirements. |  |  |
| `themeRef` _[KDexObjectReference](#kdexobjectreference)_ | themeRef is a reference to the theme that should apply to all pages bound to this host. |  | Optional: \{\} <br /> |
| `serviceAccountRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | serviceAccountRef is a reference to the service account used by the host to access secrets.<br />Each Secret must match one of the following:<br />- is annotated with 'kdex.dev/secret-type = auth-client' (multiple)<br />    - must contain key 'client-id' OR 'client_id'<br />    - may contain key 'public' (true\|false, default: false)<br />    - if not public, must contain key 'client-secret' OR 'client_secret'<br />    - must contain key 'redirect-uris' OR 'redirect_uris' (comma separated list)<br />    - may contain key 'allowed-grant-types' OR 'allowed_grant_types' (comma separated list)<br />    - may contain key 'allowed-scopes' OR 'allowed_scopes' (comma separated list)<br />    - may contain key 'require-pkce' OR 'require_pkce' (true\|false, default: false)<br />    - may contain key 'name'<br />    - may contain key 'description'<br />- is annotated with 'kdex.dev/secret-type = git' (single)<br />    - must contain key 'host'<br />    - must contain key 'org'<br />    - must contain key 'password'<br />    - must contain key 'repo'<br />    - must contain key 'username'<br />- is annotated with 'kdex.dev/secret-type = jwt-keys' (multiple)<br />    - must contain key 'private-key'<br />    - may be annotated with 'kdex.dev/active-key = true'<br />- is annotated with 'kdex.dev/secret-type = ldap' (single)<br />    - must contain key 'active-directory' (true\|false)<br />    - must contain key 'addr'<br />    - must contain key 'base-dn'<br />    - must contain key 'bind-dn'<br />    - must contain key 'bind-user'<br />    - must contain key 'bind-pass'<br />    - must contain key 'user-filter'<br />    - may contain key 'attributes' (comma separated list of attributes to retrieve)<br />- is annotated with 'kdex.dev/secret-type = npm' (single)<br />    - must contain key '.npmrc' (formatted as a complete .npmrc file)<br />- is annotated with 'kdex.dev/secret-type = oidc-client' (single)<br />    - must contain key 'client-id' OR 'client_id'<br />    - must contain key 'client-secret' OR 'client_secret'<br />    - may contain key 'block-key' OR 'block_key'<br />- is annotated with 'kdex.dev/secret-type = subject' (multiple)<br />    - must contain key 'sub'<br />    - must contain key 'password'<br />    - may contain arbitrary key(string)/value(string\|yaml) pairs which can be mapped to the claims using the spec.auth.claimMappings<br />- is of type 'kubernetes.io/dockerconfigjson' (multiple)<br />- is of type 'kubernetes.io/tls' (single) |  | Required: \{\} <br /> |
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
| `spec` _[KDexInternalPackageReferencesSpec](#kdexinternalpackagereferencesspec)_ | spec defines the desired state of KDexInternalPackageReferences |  | Required: \{\} <br /> |


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
| `builderImage` _string_ |  |  | Required: \{\} <br /> |
| `builderImagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ |  |  | Optional: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this internal package references is for. |  | Required: \{\} <br /> |
| `npmSecretRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ |  |  | Optional: \{\} <br /> |
| `packageReferences` _[PackageReference](#packagereference) array_ |  |  | MinItems: 1 <br /> |
| `serviceAccountRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ |  |  | Required: \{\} <br /> |


#### KDexInternalTranslation



KDexInternalTranslation is the Schema for the kdexinternaltranslations API



_Appears in:_
- [KDexInternalTranslationList](#kdexinternaltranslationlist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexInternalTranslation` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
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
| `contentEntries` _[ContentEntry](#contententry) array_ | contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or KDexApp references. |  | ExactlyOneOf: [appRef rawHTML] <br />MaxItems: 8 <br />MinItems: 1 <br />Required: \{\} <br /> |
| `overrideFooterRef` _[KDexObjectReference](#kdexobjectreference)_ | overrideFooterRef is an optional reference to a KDexPageFooter resource. If not specified, the footer from the archetype will be used. |  | Optional: \{\} <br /> |
| `overrideHeaderRef` _[KDexObjectReference](#kdexobjectreference)_ | overrideHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, the header from the archetype will be used. |  | Optional: \{\} <br /> |
| `overrideNavigationRefs` _object (keys:string, values:[KDexObjectReference](#kdexobjectreference))_ | overrideNavigationRefs is an optional map of keyed navigation object references. When not empty, the 'main' key must be specified. These navigations will be merged with the navigations from the archetype. |  | MaxProperties: 10 <br />Optional: \{\} <br /> |
| `pageArchetypeRef` _[KDexObjectReference](#kdexobjectreference)_ | pageArchetypeRef is a reference to the KDexPageArchetype that this binding is for. |  | Required: \{\} <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexInternalHost that this utility page belongs to. |  | Required: \{\} <br /> |




#### KDexObjectReference







_Appears in:_
- [Builder](#builder)
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
- [KDexFaaSAdaptorStatus](#kdexfaasadaptorstatus)
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
| `contentEntries` _[ContentEntry](#contententry) array_ | contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or KDexApp references. |  | ExactlyOneOf: [appRef rawHTML] <br />MaxItems: 32 <br />MinItems: 1 <br />Required: \{\} <br /> |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this binding is for. |  | Required: \{\} <br /> |
| `label` _string_ | label is the value used in menus and page titles before localization occurs (or when no translation exists for the current language). |  | MaxLength: 256 <br />MinLength: 3 <br />Required: \{\} <br /> |
| `tags` _[Tag](#tag) array_ | Tags are used for grouping and searching functions. |  | MaxItems: 16 <br />Optional: \{\} <br /> |
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


#### KDexRole



KDexRole is the Schema for the kdexroles API



_Appears in:_
- [KDexRoleList](#kdexrolelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexRole` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexRoleSpec](#kdexrolespec)_ | spec defines the desired state of KDexRole |  | Required: \{\} <br /> |


#### KDexRoleBinding



KDexRoleBinding is the Schema for the kdexrolebindings API



_Appears in:_
- [KDexRoleBindingList](#kdexrolebindinglist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexRoleBinding` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `spec` _[KDexRoleBindingSpec](#kdexrolebindingspec)_ | spec defines the desired state of KDexRoleBinding |  | Required: \{\} <br /> |


#### KDexRoleBindingList



KDexRoleBindingList contains a list of KDexRoleBinding





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexRoleBindingList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexRoleBinding](#kdexrolebinding) array_ |  |  |  |


#### KDexRoleBindingSpec



KDexRoleBindingSpec defines the desired state of KDexRoleBinding



_Appears in:_
- [KDexRoleBinding](#kdexrolebinding)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this binding is for. |  | Required: \{\} <br /> |
| `roles` _string array_ | roles is a list of KDexRole names bound to this subject. |  | MinItems: 1 <br />Required: \{\} <br /> |
| `subject` _string_ | subject is the subject identifier. It should be from the OIDC provider (e.g. Google).<br />However, if the ServiceAccount referenced by the host has secrets attached labelled with<br />"kdex.dev/secret-type=subject" then it contains a local identity managed<br />through the Secret. |  | MinLength: 5 <br />Required: \{\} <br /> |


#### KDexRoleList



KDexRoleList contains a list of KDexRole





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `kdex.dev/v1alpha1` | | |
| `kind` _string_ | `KDexRoleList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[KDexRole](#kdexrole) array_ |  |  |  |


#### KDexRoleSpec



KDexRoleSpec defines the desired state of KDexRole



_Appears in:_
- [KDexRole](#kdexrole)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `hostRef` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core)_ | hostRef is a reference to the KDexHost that this binding is for. |  | Required: \{\} <br /> |
| `rules` _[PolicyRule](#policyrule) array_ | Rules holds all the PolicyRules for this KDexRole |  | MinItems: 1 <br />Required: \{\} <br /> |


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
| `spec` _[KDexScriptLibrarySpec](#kdexscriptlibraryspec)_ | spec defines the desired state of KDexScriptLibrary |  | Required: \{\} <br /> |


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
| `scripts` _[ScriptDef](#scriptdef) array_ | scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents. |  | ExactlyOneOf: [script scriptSrc] <br />MaxItems: 8 <br />Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/-/' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/-/.+` <br /> |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is an optional list of environment variables to set in the container. |  | Optional: \{\} <br /> |
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
| `assets` _[Assets](#assets)_ | assets is a set of elements that define a portable set of design rules. |  | ExactlyOneOf: [linkHref metaId style] <br />MaxItems: 32 <br /> |
| `scriptLibraryRef` _[KDexObjectReference](#kdexobjectreference)_ | scriptLibraryRef is an optional reference to a KDexScriptLibrary resource. |  | Optional: \{\} <br /> |
| `ingressPath` _string_ | ingressPath is a prefix beginning with '/-/' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.<br />This value is determined by the implementation that embeds the Backend and cannot be changed. |  | Optional: \{\} <br />Pattern: `^/-/.+` <br /> |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is an optional list of environment variables to set in the container. |  | Optional: \{\} <br /> |
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
| `contentEntries` _[ContentEntry](#contententry) array_ | contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or KDexApp references. |  | ExactlyOneOf: [appRef rawHTML] <br />MaxItems: 8 <br />MinItems: 1 <br />Required: \{\} <br /> |
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


#### Metadata



KDexFunctionMetadata defines the metadata for the function.



_Appears in:_
- [KDexFunctionMetadata](#kdexfunctionmetadata)
- [KDexPageBindingSpec](#kdexpagebindingspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `tags` _[Tag](#tag) array_ | Tags are used for grouping and searching functions. |  | MaxItems: 16 <br />Optional: \{\} <br /> |
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
| `oidcProviderURL` _string_ | oidcProviderURL is the well known URL of the OIDC provider. |  | Required: \{\} <br /> |
| `roles` _string array_ | roles is an array of additional roles that will be requested from the provider. |  | Optional: \{\} <br /> |


#### Observer







_Appears in:_
- [KDexFaaSAdaptorSpec](#kdexfaasadaptorspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `args` _string array_ | args is an optional array of arguments that will be passed to the generator command. |  | Optional: \{\} <br /> |
| `command` _string array_ | command is an optional array that contains the code generator command and any flags necessary. |  | Optional: \{\} <br /> |
| `image` _string_ | image is the image to use for observing the function state. |  | Required: \{\} <br /> |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is the environment variables to set in the observer. |  | Optional: \{\} <br /> |
| `schedule` _string_ | schedule is the schedule in Cron format, see https://en.wikipedia.org/wiki/Cron. | */5 * * * * | Optional: \{\} <br /> |
| `serviceAccountName` _string_ | serviceAccountName is the name of the service account to use for observing the function state. |  | Optional: \{\} <br /> |


#### OpenAPI



OpenAPI holds the configuration for the host's OpenAPI support.



_Appears in:_
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `typesToInclude` _[TypeToInclude](#typetoinclude) array_ | typesToInclude specifies which route types will be outputted to the OpenAPI endpoint. | [BACKEND FUNCTION PAGE SYSTEM] | Enum: [BACKEND FUNCTION PAGE SYSTEM] <br /> |


#### PackageReference



PackageReference specifies the name and version of an NPM package. Prefereably the package should be available from
the public npm registry. If the package is not available from the public npm registry, a secretRef should be
associated with the ServiceAccount named in ServiceAccountRef to authenticate to the npm registry.
That package must contain an ES module for use in the browser.



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
- [KDexRoleSpec](#kdexrolespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `resourceNames` _string array_ | resourceNames is an optional allow list of names that the rule applies to. An empty set means the rule applies to all instances of the resources. |  | Optional: \{\} <br /> |
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
| `scheme` _string_ | scheme is the scheme to use for the host. If not specified http is assumed. | http | Optional: \{\} <br /> |
| `strategy` _[RoutingStrategy](#routingstrategy)_ | strategy is the routing strategy to use. If not specified Ingress is assumed. | Ingress | Enum: [Ingress HTTPRoute] <br />Optional: \{\} <br /> |


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


#### Runtime







_Appears in:_
- [Backend](#backend)
- [KDexAppSpec](#kdexappspec)
- [KDexHostSpec](#kdexhostspec)
- [KDexInternalHostSpec](#kdexinternalhostspec)
- [KDexScriptLibrarySpec](#kdexscriptlibraryspec)
- [KDexThemeSpec](#kdexthemespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#envvar-v1-core) array_ | env is an optional list of environment variables to set in the container. |  | Optional: \{\} <br /> |
| `replicas` _integer_ | replicas is the number of desired pods. This is a pointer to distinguish between explicit<br />zero and not specified. Defaults to 1. |  | Optional: \{\} <br /> |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources defines the compute resources required by the container.<br />More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |  | Optional: \{\} <br /> |


#### ScalingConfig



ScalingConfig defines scaling parameters.



_Appears in:_
- [Executable](#executable)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `activationScale` _integer_ | activationScale controls the minimum number of replicas that will be created<br />when the Function scales up from zero. After the Function has reached this<br />scale one time, this value is ignored. This means that the Function will<br />scale down after the activation scale is reached if the actual traffic<br />received needs a smaller scale.<br />When the Function is created, the larger of activation scale and lower<br />bound is automatically chosen as the initial target scale. | 1 | Optional: \{\} <br /> |
| `initialScale` _integer_ | initialScale controls the initial target scale a Function must reach<br />immediately after it is created before it is marked as Ready. After the<br />Function has reached this scale one time, this value is ignored. This<br />means that the Function will scale down after the initial target scale<br />is reached if the actual traffic received only needs a smaller scale.<br />When the Function is created, the larger of initialScale and<br />minScale is automatically chosen as the initial target scale. | 1 | Optional: \{\} <br /> |
| `maxScale` _integer_ | maxScale controls the maximum number of replicas that each Function<br />should have. The autoscaler will attempt to never have more than this<br />number of replicas running, or in the process of being created, at any<br />one point in time. | 0 | Optional: \{\} <br /> |
| `metric` _string_ | metric defines which metric type is watched by the Autoscaler. | concurrency | Enum: [concurrency rps] <br />Optional: \{\} <br /> |
| `minScale` _integer_ | minScale controls the minimum number of replicas that each Function<br />should have. The autoscaler will attempt to never have less than this<br />number of replicas at any one point in time. | 0 | Optional: \{\} <br /> |
| `panicThresholdPercentage` _integer_ | panicThresholdPercentage defines when the Autoscaler will move from stable<br />mode into panic mode. | 200 | Maximum: 1000 <br />Minimum: 110 <br />Optional: \{\} <br /> |
| `panicWindowPercentage` _integer_ | The panic window is defined as a percentage of the stable window to<br />assure that both are relative to each other in a working way.<br />panicWindowPercentage indicates how the window over which historical<br />data is evaluated will shrink upon entering panic mode. For example,<br />a value of 10.0 means that in panic mode the window will be 10% of the<br />stable window size. | 10 | Maximum: 100 <br />Minimum: 1 <br />Optional: \{\} <br /> |
| `scaleDownDelay` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta)_ | scaleDownDelay specifies a time window which must pass at reduced<br />concurrency before a scale-down decision is applied. This can be useful,<br />for example, to keep containers around for a configurable duration to<br />avoid a cold start penalty if new requests come in. Unlike setting a lower<br />bound, the revision will eventually be scaled down if reduced concurrency<br />is maintained for the delay period. | 0s | Optional: \{\} <br /> |
| `scaleToZeroPodRetentionPeriod` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta)_ | scaleToZeroPodRetentionPeriod determines the minimum amount of time that<br />the last pod will remain active after the Autoscaler decides to scale pods<br />to zero. This can be useful to avoid a cold start penalty if new requests<br />come in. Unlike setting a lower bound, the revision will eventually be<br />scaled down if reduced concurrency is maintained for the delay period. | 0s | Optional: \{\} <br /> |
| `stableWindow` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta)_ | stableWindow defines the sliding time window over which metrics are<br />averaged to provide the input for scaling decisions when the autoscaler<br />is not in Panic mode. | 60s | Optional: \{\} <br /> |
| `target` _integer_ | target provides the Autoscaler with a target value that it tries to<br />maintain for the configured metric. This value is metric agnostic. This<br />means the target is simply an integer value, which can be applied for any<br />metric type. |  | Optional: \{\} <br /> |
| `targetUtilizationPercentage` _integer_ | targetUtilizationPercentage specifies what percentage of the previously<br />specified target should actually be targeted by the Autoscaler. This is<br />also known as specifying the hotness at which a replica runs, which causes<br />the Autoscaler to scale up before the defined hard limit is reached. | 70 | Optional: \{\} <br /> |


#### ScriptDef





_Validation:_
- ExactlyOneOf: [script scriptSrc]

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





#### Source



Source contains source information.



_Appears in:_
- [FunctionOrigin](#functionorigin)
- [KDexFunctionStatus](#kdexfunctionstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `builder` _[Builder](#builder)_ | builder is used to build the source code into an image. |  | Optional: \{\} <br /> |
| `path` _string_ | path is the path to the source code in the repository. |  | Optional: \{\} <br /> |
| `repository` _string_ | repository is the git repository address to the source code. |  | Required: \{\} <br /> |
| `revision` _string_ | revision is the git revision (tag, branch or commit hash) to the source code. |  | Required: \{\} <br /> |




#### Tag







_Appears in:_
- [KDexFunctionMetadata](#kdexfunctionmetadata)
- [KDexPageBindingSpec](#kdexpagebindingspec)
- [Metadata](#metadata)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ |  |  | Required: \{\} <br /> |
| `description` _string_ |  |  | Optional: \{\} <br /> |
| `url` _string_ |  |  | Optional: \{\} <br /> |


#### Translation







_Appears in:_
- [KDexInternalTranslationSpec](#kdexinternaltranslationspec)
- [KDexTranslationSpec](#kdextranslationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `lang` _string_ | lang is a string containing a BCP 47 language tag that identifies the set of translations.<br />See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag. |  | Required: \{\} <br /> |
| `keysAndValues` _object (keys:string, values:string)_ | keysAndValues is a map of key=/value pairs where the key is the identifier and the value is the translation of that key in the language specified by the lang property. |  | MaxProperties: 256 <br />MinProperties: 1 <br />Required: \{\} <br /> |


#### TypeToInclude

_Underlying type:_ _string_



_Validation:_
- Enum: [BACKEND FUNCTION PAGE SYSTEM]

_Appears in:_
- [OpenAPI](#openapi)

| Field | Description |
| --- | --- |
| `BACKEND` |  |
| `FUNCTION` |  |
| `PAGE` |  |
| `SYSTEM` |  |


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


