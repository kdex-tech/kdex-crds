## CRD Relationships

See [CRD_REFERENCE.md](CRD_REFERENCE.md) for reference documentation.

```mermaid
erDiagram
    KDexFunction {
        string metadata
        string api
        string function
        string backend
    }

    KDexTranslation {
        string translations
    }

    KDexHost {
        string modulePolicy
        string assets
        string defaultLang
        string domains
        string organization
    }

    KDexPageArchetype {
        string content
    }

    KDexPageBinding {
        string label
        string path
    }

    KDexTheme {
        string assets
    }

    KDexPageFooter {
        string content
    }

    KDexApp {
        string name
        string description
        string packageReference
    }

    KDexPageHeader {
        string content
    }

    KDexPageNavigation {
        string content
    }

    KDexScriptLibrary {
        string packageReference
        string scripts
    }

    KDexPageArchetype ||--o{ KDexPageBinding : "archetype"
    KDexPageArchetype ||--o{ KDexPageFooter : "default footer"
    KDexPageArchetype ||--o{ KDexPageHeader : "default header"
    KDexPageArchetype ||--o{ KDexPageNavigation : "default navigations"
    KDexPageArchetype ||--o{ KDexScriptLibrary : "script library"

    KDexHost ||--o{ KDexPageBinding : "hosts"
    KDexHost ||--o{ KDexScriptLibrary : "script library"
    KDexHost ||--o{ KDexTheme : "theme"
    KDexHost ||--o{ KDexTranslation : "hosts"


    KDexPageBinding ||--o{ KDexApp : "uses"
    KDexPageBinding ||--o{ KDexPageBinding : "parent page"
    KDexPageBinding ||--o{ KDexPageFooter : "override footer"
    KDexPageBinding ||--o{ KDexPageHeader : "override header"
    KDexPageBinding ||--o{ KDexPageNavigation : "override navigations"
    KDexPageBinding ||--o{ KDexScriptLibrary : "script library"

    KDexPageFooter ||--o{ KDexScriptLibrary : "script library"

    KDexPageHeader ||--o{ KDexScriptLibrary : "script library"

    KDexPageNavigation ||--o{ KDexScriptLibrary : "script library"

    KDexTheme ||--o{ KDexScriptLibrary : "script library"

```