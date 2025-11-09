## CRD Relationships

See [CRD_REFERENCE.md](CRD_REFERENCE.md) for reference documentation.

```mermaid
erDiagram
    KDexApp {
        string name
        string description
        string packageReference
    }

    KDexHost ||--o{ KDexPageBinding : "hosts"
    KDexHost ||--o{ KDexScriptLibrary : "script library"
    KDexHost ||--o{ KDexTheme : "default theme"
    KDexHost ||--o{ KDexTranslation : "hosts"
    KDexHost {
        string appPolicy
        string baseMeta
        string defaultLang
        string domains
        string organization
    }

    KDexPageArchetype ||--o{ KDexPageBinding : "archetype"
    KDexPageArchetype ||--o{ KDexPageFooter : "default footer"
    KDexPageArchetype ||--o{ KDexPageHeader : "default header"
    KDexPageArchetype ||--o{ KDexPageNavigation : "default navigation"
    KDexPageArchetype ||--o{ KDexPageNavigation : "extra navigations"
    KDexPageArchetype ||--o{ KDexScriptLibrary : "script library"
    KDexPageArchetype ||--o{ KDexTheme : "override theme"
    KDexPageArchetype {
        string content
    }

    KDexPageBinding ||--o{ KDexApp : "uses"
    KDexPageBinding ||--o{ KDexPageBinding : "parent page"
    KDexPageBinding ||--o{ KDexPageFooter : "override footer"
    KDexPageBinding ||--o{ KDexPageHeader : "override header"
    KDexPageBinding ||--o{ KDexPageNavigation : "override navigation"
    KDexPageBinding ||--o{ KDexScriptLibrary : "script library"
    KDexPageBinding {
        string label
        string path
    }

    KDexPageFooter ||--o{ KDexScriptLibrary : "script library"
    KDexPageFooter {
        string content
    }

    KDexPageHeader ||--o{ KDexScriptLibrary : "script library"
    KDexPageHeader {
        string content
    }

    KDexPageNavigation ||--o{ KDexScriptLibrary : "script library"
    KDexPageNavigation {
        string content
    }

    KDexScriptLibrary {
        string packageReference
        string scripts
    }

    KDexTheme ||--o{ KDexScriptLibrary : "script library"
    KDexTheme {
        string styleItems
    }

    KDexTranslation {
        string translations
    }
```