package initialise

const kustomizeFileContent = `
"""
This file was generated by the KCL auto-gen tool. DO NOT EDIT.
Editing this file might prove futile when you re-run the KCL auto-gen generate command.
"""
import regex
import kcl_plugin.kustomize
import knit.kubernetes

build = lambda kustomization: Kustomization, resources: [{str: any}] | [kubernetes.Manifest] -> [kubernetes.Manifest] {
    kustomize.build(kustomization, resources)
}

schema Kustomization:
    r"""
    Kustomization

    Attributes
    ----------
    apiVersion : str, optional
    bases : [str], optional
        DEPRECATED. Bases are relative paths or git repository URLs specifying a directory containing a kustomization.yaml file.
    commonAnnotations : KustomizationCommonAnnotations, optional
        CommonAnnotations to add to all objects
    buildMetadata : [str], optional
        BuildMetadata is a list of strings used to toggle different build options
    commonLabels : KustomizationCommonLabels, optional
         CommonLabels to add to all objects and selectors
    configMapGenerator : [ConfigMapArgs], optional
        ConfigMapGenerator is a list of configmaps to generate from local data (one configMap per list item)
    configurations : [str], optional
        Configurations is a list of transformer configuration files
    crds : [str], optional
        Crds specifies relative paths to Custom Resource Definition files. This allows custom resources to be recognized as operands, making it possible to add them to the Resources list. CRDs themselves are not modified.
    generatorOptions : GeneratorOptions, optional
        GeneratorOptions modify behavior of all ConfigMap and Secret generators
    generators : [str], optional
        Generators is a list of files containing custom generators
    helmCharts : [HelmChart], optional
        HelmCharts is a list of helm chart configuration instances
    helmGlobals : KustomizationHelmGlobals, optional
        HelmGlobals contains helm configuration that isn't chart specific
    images : [Image], optional
        Images is a list of (image name, new name, new tag or digest) for changing image names, tags or digests. This can also be achieved with a patch, but this operator is simpler to specify.
    inventory : Inventory, optional
        Inventory appends an object that contains the record of all other objects, which can be used in apply, prune and delete
    labels : [Labels], optional
        Labels to add to all objects but not selectors
    kind : str, optional
    metadata : Metadata, optional
        Contains metadata about a Resource
    namePrefix : str, optional
        NamePrefix will prefix the names of all resources mentioned in the kustomization file including generated configmaps and secrets
    nameSuffix : str, optional
        NameSuffix will suffix the names of all resources mentioned in the kustomization file including generated configmaps and secrets
    namespace : str, optional
        Namespace to add to all objects
    replacements : [ReplacementsPath | ReplacementsInline], optional
        Substitute field(s) in N target(s) with a field from a source
    openapi : KustomizationOpenapi, optional
        OpenAPI contains information about what kubernetes schema to use
    patches : [PatchesPatchPath | PatchesInlinePatch], optional
        Apply a patch to multiple resources
    patchesJson6902 : [KustomizationPatchesJson6902Items0OneOf0 | KustomizationPatchesJson6902Items0OneOf1 | KustomizationPatchesJson6902Items0OneOf2], optional
        JSONPatches is a list of JSONPatch for applying JSON patch. See http://jsonpatch.com
    patchesStrategicMerge : [str], optional
         PatchesStrategicMerge specifies the relative path to a file containing a strategic merge patch. URLs and globs are not supported
    replicas : [Replicas], optional
        Replicas is a list of (resource name, count) for changing number of replicas for a resources. It will match any group and kind that has a matching name and that is one of: Deployment, ReplicationController, Replicaset, Statefulset.
    components : [str], optional
        Components are relative paths or git repository URLs specifying a directory containing a kustomization.yaml file of Kind Component.
    secretGenerator : [SecretArgs], optional
        SecretGenerator is a list of secrets to generate from local data (one secret per list item)
    sortOptions : KustomizationSortOptionsOneOf0 | KustomizationSortOptionsOneOf1, optional
        sortOptions is used to sort the resources kustomize outputs
    transformers : [str], optional
        Transformers is a list of files containing transformers
    validators : [str], optional
        Validators is a list of files containing validators
    vars : [Var], optional
        Allows things modified by kustomize to be injected into a container specification. A var is a name (e.g. FOO) associated with a field in a specific resource instance.  The field must contain a value of type string, and defaults to the name field of the instance
    """

    apiVersion?: str
    bases?: [str]
    commonAnnotations?: KustomizationCommonAnnotations
    buildMetadata?: [str]
    commonLabels?: KustomizationCommonLabels
    configMapGenerator?: [ConfigMapArgs]
    configurations?: [str]
    crds?: [str]
    generatorOptions?: GeneratorOptions
    generators?: [str]
    helmCharts?: [HelmChart]
    helmGlobals?: KustomizationHelmGlobals
    images?: [Image]
    inventory?: Inventory
    labels?: [Labels]
    kind?: str
    metadata?: Metadata
    namePrefix?: str
    nameSuffix?: str
    namespace?: str
    replacements?: [ReplacementsPath | ReplacementsInline]
    openapi?: KustomizationOpenapi
    patches?: [PatchesPatchPath | PatchesInlinePatch]
    patchesJson6902?: [KustomizationPatchesJson6902Items0OneOf0 | KustomizationPatchesJson6902Items0OneOf1 | KustomizationPatchesJson6902Items0OneOf2]
    patchesStrategicMerge?: [str]
    replicas?: [Replicas]
    components?: [str]
    secretGenerator?: [SecretArgs]
    sortOptions?: KustomizationSortOptionsOneOf0 | KustomizationSortOptionsOneOf1
    transformers?: [str]
    validators?: [str]
    vars?: [Var]

schema ConfigMapArgs:
    r"""
    ConfigMapArgs contains the metadata of how to generate a configmap

    Attributes
    ----------
    KVSources : [Kvsource], optional
    behavior : "create" | "replace" | "merge", optional
    env : str, optional
        Deprecated.  Use envs instead.
    envs : [str], optional
        A list of file paths. The contents of each file should be one key=value pair per line
    files : [str], optional
        A list of file sources to use in creating a list of key, value pairs
    literals : [str], optional
        A list of literal pair sources. Each literal source should be a key and literal value, e.g. key=value
    name : str, optional
        Name - actually the partial name - of the generated resource
    namespace : str, optional
        Namespace for the configmap, optional
    options : GeneratorOptions, optional
        GeneratorOptions modify behavior of all ConfigMap and Secret generators
    """

    KVSources?: [Kvsource]
    behavior?: "create" | "replace" | "merge"
    env?: str
    envs?: [str]
    files?: [str]
    literals?: [str]
    name?: str
    namespace?: str
    options?: GeneratorOptions

schema FieldSelector:
    r"""
    Refers to the field of the object referred to by objref whose value will be extracted for use in replacing $(FOO)

    Attributes
    ----------
    fieldpath : str, optional
    """

    fieldpath?: str

schema FieldSpec:
    r"""
    FieldSpec

    Attributes
    ----------
    create : bool, optional
    group : str, optional
    kind : str, optional
    path : str, optional
    version : str, optional
    """

    create?: bool
    group?: str
    kind?: str
    path?: str
    version?: str

schema GeneratorOptions:
    r"""
    GeneratorOptions modify behavior of all ConfigMap and Secret generators

    Attributes
    ----------
    annotations : KustomizationSecretGeneratorItems0OptionsAnnotations, optional
        Annotations to add to all generated resources
    disableNameSuffixHash : bool, optional
        DisableNameSuffixHash if true disables the default behavior of adding a suffix to the names of generated resources that is a hash of the resource contents
    immutable : bool, optional
        Immutable if true add to all generated resources
    labels : KustomizationSecretGeneratorItems0OptionsLabels, optional
        Labels to add to all generated resources
    """

    annotations?: KustomizationSecretGeneratorItems0OptionsAnnotations
    disableNameSuffixHash?: bool
    immutable?: bool
    labels?: KustomizationSecretGeneratorItems0OptionsLabels

schema HelmChart:
    r"""
    HelmChart

    Attributes
    ----------
    name : str, optional
    version : str, optional
    repo : str, optional
    releaseName : str, optional
    namespace : str, optional
    valuesFile : str, optional
    valuesInline : KustomizationHelmChartsItems0ValuesInline, optional
    valuesMerge : "merge" | "override" | "replace", optional
    includeCRDs : bool, optional
    skipHooks : bool, optional
    additionalValuesFiles : [str], optional
    skipTests : bool, optional
    apiVersions : [str], optional
    kubeVersion : str, optional
    nameTemplate : str, optional
    """

    name?: str
    version?: str
    repo?: str
    releaseName?: str
    namespace?: str
    valuesFile?: str
    valuesInline?: KustomizationHelmChartsItems0ValuesInline
    valuesMerge?: "merge" | "override" | "replace"
    includeCRDs?: bool
    skipHooks?: bool
    additionalValuesFiles?: [str]
    skipTests?: bool
    apiVersions?: [str]
    kubeVersion?: str
    nameTemplate?: str

schema Image:
    r"""
    Image

    Attributes
    ----------
    digest : str, optional
    name : str, optional
    newName : str, optional
    newTag : str, optional
    """

    digest?: str
    name?: str
    newName?: str
    newTag?: str

schema Inventory:
    r"""
    Inventory appends an object that contains the record of all other objects, which can be used in apply, prune and delete

    Attributes
    ----------
    configMap : NameArgs, optional
    $type : str, optional
    """

    configMap?: NameArgs
    $type?: str

schema KustomizationCommonAnnotations:
    r"""
    CommonAnnotations to add to all objects
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationCommonLabels:
    r"""
     CommonLabels to add to all objects and selectors
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationConfigMapGeneratorItems0OptionsAnnotations:
    r"""
    Annotations to add to all generated resources
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationConfigMapGeneratorItems0OptionsLabels:
    r"""
    Labels to add to all generated resources
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationGeneratorOptionsAnnotations:
    r"""
    Annotations to add to all generated resources
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationGeneratorOptionsLabels:
    r"""
    Labels to add to all generated resources
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationHelmChartsItems0ValuesInline:
    r"""
    KustomizationHelmChartsItems0ValuesInline
    """

    [key: str]: any

    check:
        regex.match(key, r".*")

schema KustomizationHelmGlobals:
    r"""
    HelmGlobals contains helm configuration that isn't chart specific

    Attributes
    ----------
    chartHome : str, optional
        ChartHome is a file path, relative to the kustomization root, to a directory containing a subdirectory for each chart to be included in the kustomization
    configHome : str, optional
        ConfigHome defines a value that kustomize should pass to helm via the HELM_CONFIG_HOME environment variable
    """

    chartHome?: str
    configHome?: str

schema KustomizationLabelsItems0Pairs:
    r"""
    Pairs contains the key-value pairs for labels to add
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationMetadataAnnotations:
    r"""
    KustomizationMetadataAnnotations
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationMetadataLabels:
    r"""
    KustomizationMetadataLabels
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationOpenapi:
    r"""
    OpenAPI contains information about what kubernetes schema to use
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationPatchesJson6902Items0OneOf0:
    r"""
    KustomizationPatchesJson6902Items0OneOf0

    Attributes
    ----------
    path : str, required
        relative file path for a json patch file inside a kustomization
    target : PatchTarget, required
        Refers to a Kubernetes object that the json patch will be applied to. It must refer to a Kubernetes resource under the purview of this kustomization
    """

    path: str
    target: PatchTarget

schema KustomizationPatchesJson6902Items0OneOf1:
    r"""
    KustomizationPatchesJson6902Items0OneOf1

    Attributes
    ----------
    patch : str, required
        inline json patch
    target : PatchTarget, required
        Refers to a Kubernetes object that the json patch will be applied to. It must refer to a Kubernetes resource under the purview of this kustomization
    """

    patch: str
    target: PatchTarget

schema KustomizationPatchesJson6902Items0OneOf2:
    r"""
    KustomizationPatchesJson6902Items0OneOf2

    Attributes
    ----------
    op : "add" | "remove" | "replace" | "move" | "copy" | "test", required
        The operation
    from : str, optional
        The source location.
    path : str, required
        The target location.
    value : str | [str], optional
    """

    op: "add" | "remove" | "replace" | "move" | "copy" | "test"
    from?: str
    path: str
    value?: str | [str]

schema KustomizationReplacementsItems0OneOf1Source:
    r"""
    The source of the value

    Attributes
    ----------
    group : str, optional
        The group of the referent
    version : str, optional
        The version of the referent
    kind : str, optional
        The kind of the referent
    name : str, optional
        The name of the referent
    namespace : str, optional
        The namespace of the referent
    fieldPath : str, optional
        The structured path to the source value
    options : KustomizationReplacementsItems0OneOf1SourceOptions, optional
    """

    group?: str
    version?: str
    kind?: str
    name?: str
    namespace?: str
    fieldPath?: str
    options?: KustomizationReplacementsItems0OneOf1SourceOptions

schema KustomizationReplacementsItems0OneOf1SourceOptions:
    r"""
    KustomizationReplacementsItems0OneOf1SourceOptions

    Attributes
    ----------
    delimiter : str, optional
    index : float, optional
    create : bool, optional
    """

    delimiter?: str
    index?: float
    create?: bool

schema KustomizationReplacementsItems0OneOf1TargetsItems0:
    r"""
    KustomizationReplacementsItems0OneOf1TargetsItems0

    Attributes
    ----------
    select : KustomizationReplacementsItems0OneOf1TargetsItems0Select, required
        Include objects that match this
    reject : [KustomizationReplacementsItems0OneOf1TargetsItems0RejectItems0], optional
        Exclude objects that match this
    fieldPaths : [str], optional
        The structured path(s) to the target nodes
    options : KustomizationReplacementsItems0OneOf1TargetsItems0Options, optional
    """

    select: KustomizationReplacementsItems0OneOf1TargetsItems0Select
    reject?: [KustomizationReplacementsItems0OneOf1TargetsItems0RejectItems0]
    fieldPaths?: [str]
    options?: KustomizationReplacementsItems0OneOf1TargetsItems0Options

schema KustomizationReplacementsItems0OneOf1TargetsItems0Options:
    r"""
    KustomizationReplacementsItems0OneOf1TargetsItems0Options

    Attributes
    ----------
    delimiter : str, optional
    index : float, optional
    create : bool, optional
    """

    delimiter?: str
    index?: float
    create?: bool

schema KustomizationReplacementsItems0OneOf1TargetsItems0RejectItems0:
    r"""
    Exclude objects that match this

    Attributes
    ----------
    group : str, optional
        The group of the referent
    version : str, optional
        The version of the referent
    kind : str, optional
        The kind of the referent
    name : str, optional
        The name of the referent
    namespace : str, optional
        The namespace of the referent
    """

    group?: str
    version?: str
    kind?: str
    name?: str
    namespace?: str

schema KustomizationReplacementsItems0OneOf1TargetsItems0Select:
    r"""
    Include objects that match this

    Attributes
    ----------
    group : str, optional
        The group of the referent
    version : str, optional
        The version of the referent
    kind : str, optional
        The kind of the referent
    name : str, optional
        The name of the referent
    namespace : str, optional
        The namespace of the referent
    """

    group?: str
    version?: str
    kind?: str
    name?: str
    namespace?: str

schema KustomizationSecretGeneratorItems0OptionsAnnotations:
    r"""
    Annotations to add to all generated resources
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationSecretGeneratorItems0OptionsLabels:
    r"""
    Labels to add to all generated resources
    """

    [key: str]: str

    check:
        regex.match(key, r".*")

schema KustomizationSortOptionsOneOf0:
    r"""
    KustomizationSortOptionsOneOf0

    Attributes
    ----------
    order : "legacy", optional
    legacySortOptions : KustomizationSortOptionsOneOf0LegacySortOptions, optional
    """

    order?: "legacy"
    legacySortOptions?: KustomizationSortOptionsOneOf0LegacySortOptions

schema KustomizationSortOptionsOneOf0LegacySortOptions:
    r"""
    KustomizationSortOptionsOneOf0LegacySortOptions

    Attributes
    ----------
    orderFirst : any, optional
    orderLast : any, optional
    """

    orderFirst?: any
    orderLast?: any

schema KustomizationSortOptionsOneOf1:
    r"""
    KustomizationSortOptionsOneOf1

    Attributes
    ----------
    order : "fifo", optional
    """

    order?: "fifo"

schema Kvsource:
    r"""
    Kvsource

    Attributes
    ----------
    args : [str], optional
    name : str, optional
    pluginType : str, optional
    """

    args?: [str]
    name?: str
    pluginType?: str

schema Labels:
    r"""
    Labels

    Attributes
    ----------
    pairs : KustomizationLabelsItems0Pairs, optional
        Pairs contains the key-value pairs for labels to add
    includeSelectors : bool, optional
        IncludeSelectors indicates should transformer include the fieldSpecs for selectors
    includeTemplates : bool, optional
        IncludeTemplates indicates should transformer include the template labels
    fields : [FieldSpec], optional
        FieldSpec completely specifies a kustomizable field in a k8s API object. It helps define the operands of transformations
    """

    pairs?: KustomizationLabelsItems0Pairs
    includeSelectors?: bool
    includeTemplates?: bool
    fields?: [FieldSpec]

schema Metadata:
    r"""
    Contains metadata about a Resource

    Attributes
    ----------
    name : str, optional
    namespace : str, optional
    labels : KustomizationMetadataLabels, optional
    annotations : KustomizationMetadataAnnotations, optional
    """

    name?: str
    namespace?: str
    labels?: KustomizationMetadataLabels
    annotations?: KustomizationMetadataAnnotations

schema NameArgs:
    r"""
    NameArgs

    Attributes
    ----------
    name : str, optional
    namespace : str, optional
    """

    name?: str
    namespace?: str

schema PatchTarget:
    r"""
    Refers to a Kubernetes object that the json patch will be applied to. It must refer to a Kubernetes resource under the purview of this kustomization

    Attributes
    ----------
    group : str, optional
    kind : str, required
    name : str, required
    namespace : str, optional
    version : str, required
    """

    group?: str
    kind: str
    name: str
    namespace?: str
    version: str

schema PatchTargetOptional:
    r"""
    Refers to a Kubernetes object that the patch will be applied to. It must refer to a Kubernetes resource under the purview of this kustomization

    Attributes
    ----------
    group : str, optional
    kind : str, optional
    name : str, optional
    namespace : str, optional
    version : str, optional
    labelSelector : str, optional
    annotationSelector : str, optional
    """

    group?: str
    kind?: str
    name?: str
    namespace?: str
    version?: str
    labelSelector?: str
    annotationSelector?: str

schema PatchesInlinePatch:
    r"""
    PatchesInlinePatch

    Attributes
    ----------
    options : PatchesOptions, optional
    patch : str, required
    target : PatchTargetOptional, optional
        Refers to a Kubernetes object that the patch will be applied to. It must refer to a Kubernetes resource under the purview of this kustomization
    """

    options?: PatchesOptions
    patch: str
    target?: PatchTargetOptional

schema PatchesOptions:
    r"""
    PatchesOptions

    Attributes
    ----------
    allowNameChange : bool, optional
    allowKindChange : bool, optional
    """

    allowNameChange?: bool
    allowKindChange?: bool

schema PatchesPatchPath:
    r"""
    PatchesPatchPath

    Attributes
    ----------
    options : PatchesOptions, optional
    path : str, required
    target : PatchTargetOptional, optional
        Refers to a Kubernetes object that the patch will be applied to. It must refer to a Kubernetes resource under the purview of this kustomization
    """

    options?: PatchesOptions
    path: str
    target?: PatchTargetOptional

schema ReplacementsInline:
    r"""
    ReplacementsInline

    Attributes
    ----------
    source : KustomizationReplacementsItems0OneOf1Source, required
        The source of the value
    targets : [KustomizationReplacementsItems0OneOf1TargetsItems0], required
        The N fields to write the value to
    """

    source: KustomizationReplacementsItems0OneOf1Source
    targets: [KustomizationReplacementsItems0OneOf1TargetsItems0]

schema ReplacementsPath:
    r"""
    ReplacementsPath

    Attributes
    ----------
    path : str, required
    """

    path: str

schema Replicas:
    r"""
    Replicas

    Attributes
    ----------
    name : str, optional
    count : float, optional
    """

    name?: str
    count?: float

schema SecretArgs:
    r"""
    SecretArgs contains the metadata of how to generate a secret

    Attributes
    ----------
    KVSources : [Kvsource], optional
    behavior : "create" | "replace" | "merge", optional
    env : str, optional
    envs : [str], optional
    files : [str], optional
    literals : [str], optional
    name : str, optional
        Name - actually the partial name - of the generated resource
    namespace : str, optional
        Namespace for the secret, optional
    options : GeneratorOptions, optional
        GeneratorOptions modify behavior of all ConfigMap and Secret generators
    $type : str, optional
        Type of the secret, optional
    """

    KVSources?: [Kvsource]
    behavior?: "create" | "replace" | "merge"
    env?: str
    envs?: [str]
    files?: [str]
    literals?: [str]
    name?: str
    namespace?: str
    options?: GeneratorOptions
    $type?: str

schema Target:
    r"""
    Refers to a Kubernetes resource under the purview of this kustomization

    Attributes
    ----------
    apiVersion : str, optional
    group : str, optional
    kind : str, optional
    name : str, required
    version : str, optional
    """

    apiVersion?: str
    group?: str
    kind?: str
    name: str
    version?: str

schema Var:
    r"""
    Represents a variable whose value will be sourced from a field in a Kubernetes object.

    Attributes
    ----------
    fieldref : FieldSelector, optional
        Refers to the field of the object referred to by objref whose value will be extracted for use in replacing $(FOO)
    name : str, required
        Value of identifier name e.g. FOO used in container args, annotations, Appears in pod template as $(FOO)
    objref : Target, required
        Refers to a Kubernetes resource under the purview of this kustomization
    """

    fieldref?: FieldSelector
    name: str
    objref: Target


`