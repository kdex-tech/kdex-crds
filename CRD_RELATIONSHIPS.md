## CRD Relationships

See [CRD_REFERENCE.md](CRD_REFERENCE.md) for reference documentation.

```mermaid
erDiagram
    KDexHost ||--o{ KDexPageBinding : "hosts"
    KDexHost ||--o{ KDexRenderPage : "hosts"
    KDexHost ||--o{ KDexTranslation : "hosts"
    KDexHost ||--o{ KDexTheme : "default theme"

    KDexPageArchetype ||--o{ KDexPageBinding : "archetype"
    KDexPageArchetype ||--o{ KDexPageFooter : "default footer"
    KDexPageArchetype ||--o{ KDexPageHeader : "default header"
    KDexPageArchetype ||--o{ KDexPageNavigation : "default navigation"
    KDexPageArchetype ||--o{ KDexPageNavigation : "extra navigations"
    KDexPageArchetype ||--o{ KDexTheme : "override theme"

    KDexPageBinding ||--o{ KDexApp : "uses"
    KDexPageBinding ||--o{ KDexPageFooter : "override footer"
    KDexPageBinding ||--o{ KDexPageHeader : "override header"
    KDexPageBinding ||--o{ KDexPageNavigation : "override navigation"
    KDexPageBinding ||--o{ KDexPageBinding : "parent page"

    KDexRenderPage ||--o{ KDexRenderPage : "parent page"
    KDexRenderPage ||--o{ KDexTheme : "theme"

    KDexApp {
        string name
        string description
        string packageReference
    }

    KDexHost {
        string appPolicy
        string baseMeta
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

    KDexPageFooter {
        string content
    }

    KDexPageHeader {
        string content
    }

    KDexPageNavigation {
        string content
    }

    KDexRenderPage {
        string path
    }

    KDexTheme {
        string styleItems
    }

    KDexTranslation {
        string translations
    }
```