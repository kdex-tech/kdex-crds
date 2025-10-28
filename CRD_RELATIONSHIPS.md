## CRD Relationships

See [CRD_REFERENCE.md](CRD_REFERENCE.md) for reference documentation.

```mermaid
erDiagram
    MicroFrontEndHost ||--o{ MicroFrontEndPageBinding : "hosts"
    MicroFrontEndHost ||--o{ MicroFrontEndRenderPage : "hosts"
    MicroFrontEndHost ||--o{ MicroFrontEndTranslation : "hosts"
    MicroFrontEndHost ||--o{ MicroFrontendTheme : "default theme"

    MicroFrontEndPageArchetype ||--o{ MicroFrontEndPageBinding : "archetype"
    MicroFrontEndPageArchetype ||--o{ MicroFrontEndPageFooter : "default footer"
    MicroFrontEndPageArchetype ||--o{ MicroFrontEndPageHeader : "default header"
    MicroFrontEndPageArchetype ||--o{ MicroFrontEndPageNavigation : "default navigation"
    MicroFrontEndPageArchetype ||--o{ MicroFrontendTheme : "override theme"

    MicroFrontEndPageBinding ||--o{ MicroFrontEndApp : "uses"
    MicroFrontEndPageBinding ||--o{ MicroFrontEndPageFooter : "override footer"
    MicroFrontEndPageBinding ||--o{ MicroFrontEndPageHeader : "override header"
    MicroFrontEndPageBinding ||--o{ MicroFrontEndPageNavigation : "override navigation"
    MicroFrontEndPageBinding ||--o{ MicroFrontEndPageBinding : "parent page"

    MicroFrontEndRenderPage ||--o{ MicroFrontEndRenderPage : "parent page"

    MicroFrontEndApp {
        string name
        string description
        string packageReference
    }

    MicroFrontEndHost {
        string appPolicy
        string baseMeta
        string defaultLang
        string domains
        string organization
    }

    MicroFrontEndPageArchetype {
        string content
    }

    MicroFrontEndPageBinding {
        string label
        string path
    }

    MicroFrontEndPageFooter {
        string content
    }

    MicroFrontEndPageHeader {
        string content
    }

    MicroFrontEndPageNavigation {
        string content
    }

    MicroFrontEndRenderPage {
        string path
    }

    MicroFrontendTheme {
        string styleItems
    }

    MicroFrontEndTranslation {
        string translations
    }
```