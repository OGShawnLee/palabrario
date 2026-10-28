# 🎨 Palabrario — Sistema de Diseño, Paleta de Color y Reglas Visuales (`.context/PALETTE.md`)

> **Propósito de este documento (Design System & Visual SSOT):**  
> Este documento constituye la **Fuente Única de Verdad (SSOT) Visual** de **Palabrario** (`Vue 3 + TypeScript + UnoCSS + Reka UI`), compartida por `local-app` (Wails V3 + Capacitor PWA) y el futuro `web-portal` (Astro) a través de `@palabrario/ui`.  
> Define una **Paleta Semántica Podada** extraída de **una única ilustración de referencia**, evitando la dispersión de colores y garantizando refactorización instantánea desde un único archivo: [`packages/ui/src/uno.config.ts`](../packages/ui/src/uno.config.ts).

---

## 1. Filosofía Estética e Ilustración de Referencia: *"Imprenta Táctil Ilustrada"*

La identidad visual de **Palabrario** nace de la **ilustración editorial amigable**: lienzo crema mantequilla, tinta grafito suave (nunca negro duro), un acento mandarina pastel protagónico y pequeños acentos pastel (lavanda, salvia, miel). La interfaz debe sentirse **cálida, ligera y acogedora**, invitando al descubrimiento sin intimidar al estudiante.

| Ilustración Canónica Única |
| :---: |
| ![Ilustración Amigable y Táctil](./assets/references/ref-tactile-friendly.png) |
| *Aporta todo el sistema: fondo crema cálido, tinta grafito suave, mandarina del corazón y el florero, lavanda de la corona, salvia del pasto y miel cálida. Formas muy redondeadas y profundidad suave sin contornos duros.* |

**Principios derivados:**
1. **Suavidad sobre contraste duro:** Profundidad mediante **sombras difusas de baja opacidad** teñidas con la tinta (`shadow-soft`), nunca sombras sólidas desplazadas ni bordes gruesos de tinta.
2. **Jerarquía por opacidad:** Bordes, divisores y textos secundarios se resuelven con modificadores de opacidad del token `ink` (`border-ink/10`, `text-ink/65`), no con nuevos colores.
3. **Pasteles por transparencia:** Los fondos pastel de badges y avisos se logran con `/10–/20` sobre los tokens de acento (`bg-brand/10`, `bg-success/15`).

---

## 2. Arquitectura de Color: Tokens Semánticos Podados (Fácil Refactorización)

> [!IMPORTANT]
> **REGLA DE ORO CONTRA LA INFLACIÓN DE COLORES:**  
> Queda **estrictamente prohibido** usar valores hexadecimales sueltos (`bg-[#...]`) o colores por defecto de UnoCSS/Tailwind en los componentes `.vue`.  
> Todos los componentes utilizan exclusivamente **4 Tokens Estructurales Base** y **3 Tokens de Acento Semántico**. Las variaciones de intensidad se logran con **modificadores de opacidad nativos de UnoCSS** (`/4`, `/10`, `/15`, `/20`, `/65`), permitiendo cambiar la apariencia de toda la aplicación editando únicamente `uno.config.ts`.

### 2.1. Los 4 Pilares Estructurales (El 90% de la Interfaz)

| Token Semántico | Valor Hex | Origen en la Ilustración | Rol Arquitectónico y Reglas de Uso | Variaciones por Opacidad |
| :--- | :---: | :--- | :--- | :--- |
| **`paper`** | `#FFFDF5` | Crema mantequilla del lienzo | **Lienzo Principal (`body` / `AppLayout`):** Fondo general cálido que evita el blanco clínico. | `bg-paper/85` + `backdrop-blur-md` para `TopBar` flotante. |
| **`surface`** | `#FFFFFF` | Blanco limpio de los pétalos | **Superficie Elevada:** Tarjetas (`WordCard`, `DeckCard`), modales (`Reka Dialog`), inputs, botones secundarios y pestañas activas. | `bg-surface` puro para destacar sobre `bg-paper`. |
| **`ink`** | `#3B3B39` | Carbón grafito de los trazos | **Tinta Universal (Texto, Iconos, Bordes y Sombras):** Sustituye al negro puro en todo el sistema. | • `text-ink`: Texto primario<br>• `text-ink/65`–`/75`: Texto secundario / glosas<br>• `text-ink/40`–`/45`: Placeholders, metadatos y *eyebrows*<br>• `border-ink/8`–`/12`: Bordes y divisores<br>• `bg-ink/4`–`/6`: Fondos sutiles (`card-subtle`, píldoras neutras, hover) |
| **`brand`** | `#F78C56` | Mandarina pastel del corazón y el florero | **Acción Primaria:** CTAs (`btn-primary`), racha (`🔥`), temporizador de Arena, logo y acentos hero. | • `bg-brand`: Botón primario<br>• `bg-brand-soft` (`#FCA374`): Hover luminoso<br>• `bg-brand/10`–`/15` + `border-brand/25`: *Friendly blurb* y badges cálidos<br>• `ring-brand/15`: Anillo de foco |

### 2.2. Los 3 Acentos Pastel Semánticos (El 10% Restante — Solo Estados y Categorías)

| Token Semántico | Valor Hex | Origen en la Ilustración | Uso Exclusivo Permitido en la UI | Acabado Pastel Canónico |
| :--- | :---: | :--- | :--- | :--- |
| **`info`** | `#7C92FF` | Azul lavanda del estambre y la corona | **Regiones, XP y Nivel CEFR:** Etiquetas dialectales, contador de XP, estado `NEW` en Studio, citas de ejemplo (borde de `blockquote`). | `bg-info/15 text-info border border-info/25` |
| **`success`** | `#6CBF73` | Verde salvia del pasto | **Aprobación y Acierto:** Sentidos `APPROVED` por Quórum (**RN06**), respuestas correctas en Arena (**CU15**) y nivel `MASTERED` en SM-2 (**CU13**). | `bg-success/12`–`/15 text-success border-success/50`; sólido `bg-success text-surface` para el indicador seleccionado. |
| **`warning`** | `#F2B84B` | Miel cálido | **Atención, Quórum en Curso y Tabú:** Propuestas en `PENDING_QUORUM`, advertencias de tabú/vulgaridad, medallas de podio y estado `LEARNING`. | `bg-warning/20 text-ink border border-warning/40` *(texto en `ink` por contraste).* |

---

## 3. Tipografía Unificada: Familia `Satoshi`

Toda la interfaz utiliza **`Satoshi`** como familia única, diferenciando la jerarquía exclusivamente mediante peso, tamaño y espaciado:

| Rol Tipográfico | Clases UnoCSS Canónicas | Aplicación en Palabrario |
| :--- | :--- | :--- |
| **Display / Lema Hero** | `font-black text-3xl md:text-4xl tracking-tight text-ink` | Lema en `WordDetail`, Palabra del Día y encabezados de Arena. |
| **Títulos de Sección (`H1–H2`)** | `font-bold text-xl md:text-2xl tracking-tight text-ink` | Encabezados de vistas (`Studio`, `Diccionario`, `Cola de Quórum`). |
| **Subtítulos y Tarjetas (`H3`)** | `font-bold text-base md:text-lg text-ink` | Títulos de `DeckCard`, `CollectionCard` y nombres de usuario. |
| **Cuerpo de Lectura y Definiciones** | `text-base leading-relaxed text-ink/85` | Definiciones, *friendly blurbs* y ejemplos de uso. |
| **Eyebrows / Rótulos de Sección** | `text-xs font-bold uppercase tracking-wider text-ink/45` | *"Selecciona la definición correcta"*, *"¿Cómo se dice en otras regiones?"*. |
| **Badges y Botones** | `text-xs font-bold tracking-wide` (badges) · `font-bold` / `font-semibold` (botones) | Badges de país, CEFR, estados de quórum y acciones. |
| **Códigos PIN, Reloj y Puntajes** | `font-black tabular-nums` | PIN de Arena (`CU14`), cuenta regresiva (`RN10`), racha y XP. |

---

## 4. Geometría de Elementos y Componentes ("Soft Tactile")

### 4.1. Esquinas, Bordes y Elevación
1. **Esquinas Muy Amigables:**
   - Contenedores principales y tarjetas (`card-editorial`): `rounded-3xl`.
   - Botones, inputs, opciones de Arena, avisos internos y `card-subtle`: `rounded-2xl`.
   - Pestañas internas y botones fantasma: `rounded-xl`.
   - Badges, píldoras y contadores: `rounded-full`.
2. **Bordes por Opacidad (Sin Trazos Duros):**
   - Tarjetas y superficies: `border border-ink/10`. Inputs y botones secundarios: `border border-ink/12`.
   - Divisores internos: `border-ink/8`. **Prohibido** `border-2 border-ink` o bordes de tinta al 100%.
3. **Elevación Suave (`shadow-soft`):**
   - **Reposo:** `shadow-soft` (tarjetas, inputs, botones secundarios, pestaña activa).
   - **Elevado / Modales:** `shadow-soft-lg`.
   - **CTA de Marca:** `shadow-soft-brand` (halo mandarina en `btn-primary`).
   - **Hover interactivo:** `-translate-y-0.5`; **Active:** `translate-y-0`. Transiciones `transition-all duration-200`.

### 4.2. Anatomía de los Elementos Clave
- **Botón Primario (`btn-primary`):** Mandarina sólida, texto `surface`, `rounded-2xl`, `shadow-soft-brand`, hover a `brand-soft` con leve elevación.
- **Botón Secundario (`btn-secondary`):** `bg-surface` con `border-ink/12` y `shadow-soft`; hover `bg-ink/4`.
- **Botón Fantasma (`btn-ghost`):** Transparente, `text-ink/70`, hover `bg-ink/6`. Para acciones terciarias y menús.
- **Tarjeta Principal (`card-editorial`):** `bg-surface rounded-3xl border-ink/10 shadow-soft p-6`.
- **Tarjeta Sutil (`card-subtle`):** `bg-ink/4 rounded-2xl border-ink/8 p-4`. Para bloques informativos anidados (heterónimos, metadatos).
- **Navegación por Píldoras (TopBar / Tabs):** Contenedor `bg-ink/4 p-1 rounded-2xl border border-ink/6`; pestaña activa `bg-surface text-ink shadow-soft rounded-xl`; inactiva `text-ink/55 hover:text-ink`.
- **`FriendlyBlurb` (Estilo Vocabulary.com):** Aviso cálido `bg-brand/10 border border-brand/25 rounded-2xl p-5` con eyebrow `text-brand` y cuerpo `text-ink/85 leading-relaxed`; párrafo de etimología en `text-ink/65`.
- **Contenedor de `SenseItem`:** `card-editorial` con listón superior `badge-pill bg-brand text-surface` (*"Acepciones de {lemma}"*), índice numérico en círculo `bg-ink/6 text-ink/70`, definición `font-semibold`, cita de ejemplo en `blockquote` con `border-l-2 border-info/50 text-ink/60 italic`, píldoras de región/CEFR y desplegable `"Ver más ⌄"` (**Reka UI `Collapsible`**).
- **Opciones de Arena:** `rounded-2xl border p-3.5`; reposo `bg-surface border-ink/25 text-ink/75`; seleccionada `bg-success/12 border-success/50 shadow-soft` con índice sólido `bg-success text-surface`.
- **Inputs y Omnibox (`input-editorial`):** `bg-surface rounded-2xl border-ink/12 shadow-soft px-4 py-3`, foco `border-brand ring-3 ring-brand/15`.

---

## 5. Configuración Canónica de UnoCSS (`packages/ui/src/uno.config.ts`)

Esta configuración vive **únicamente** en `@palabrario/ui` y es re-exportada por cada shell (`local-app/frontend/uno.config.ts` → `export { default } from '@palabrario/ui/uno.config'`). Incluye `content.filesystem` para garantizar la extracción de clases de los componentes compartidos desde cualquier workspace consumidor.

```typescript
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig, presetWind3, presetIcons } from 'unocss'

const uiRoot = path.dirname(fileURLToPath(import.meta.url))
const sharedComponentsGlob = path.resolve(uiRoot, 'components/**/*.{vue,ts}').replace(/\\/g, '/')

export default defineConfig({
  content: {
    filesystem: [sharedComponentsGlob],
  },
  presets: [
    presetWind3(),
    presetIcons({
      scale: 1.2,
      extraProperties: {
        'display': 'inline-block',
        'vertical-align': 'middle',
      },
    }),
  ],
  theme: {
    colors: {
      // 1. Los 4 Pilares Estructurales
      paper: '#FFFDF5',        // Crema mantequilla suave del lienzo
      surface: '#FFFFFF',      // Superficie limpia para tarjetas elevadas
      ink: '#3B3B39',          // Carbón grafito suave (nunca negro duro)
      brand: {
        DEFAULT: '#F78C56',    // Mandarina pastel del corazón y florero
        soft: '#FCA374',       // Estado hover luminoso
      },
      // 2. Los 3 Acentos Pastel
      info: '#7C92FF',         // Azul lavanda
      success: '#6CBF73',      // Verde salvia
      warning: '#F2B84B',      // Miel cálido
    },
    fontFamily: {
      sans: '"Satoshi", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    },
    boxShadow: {
      'soft': '0 4px 20px -2px rgba(59, 59, 57, 0.06), 0 2px 6px -1px rgba(59, 59, 57, 0.04)',
      'soft-lg': '0 12px 32px -4px rgba(59, 59, 57, 0.10), 0 4px 12px -2px rgba(59, 59, 57, 0.05)',
      'soft-brand': '0 8px 20px -4px rgba(247, 140, 86, 0.35)',
    },
  },
  shortcuts: {
    'card-editorial': 'bg-surface text-ink rounded-3xl border border-ink/10 shadow-soft p-6 transition-all duration-200',
    'card-subtle': 'bg-ink/4 text-ink rounded-2xl border border-ink/8 p-4',
    'btn-primary': 'inline-flex items-center justify-center gap-2 bg-brand text-surface font-bold px-5 py-2.5 rounded-2xl shadow-soft-brand hover:(bg-brand-soft -translate-y-0.5) active:translate-y-0 transition-all duration-200 cursor-pointer disabled:(opacity-50 pointer-events-none)',
    'btn-secondary': 'inline-flex items-center justify-center gap-2 bg-surface text-ink font-semibold px-5 py-2.5 rounded-2xl border border-ink/12 shadow-soft hover:(bg-ink/4 border-ink/20 -translate-y-0.5) active:translate-y-0 transition-all duration-200 cursor-pointer disabled:(opacity-50 pointer-events-none)',
    'btn-ghost': 'inline-flex items-center justify-center gap-2 bg-transparent text-ink/70 font-semibold px-4 py-2 rounded-xl hover:(bg-ink/6 text-ink) transition-all cursor-pointer',
    'input-editorial': 'w-full bg-surface text-ink placeholder:text-ink/40 rounded-2xl border border-ink/12 px-4 py-3 shadow-soft focus:(outline-none border-brand ring-3 ring-brand/15) transition-all',
    'badge-pill': 'inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold tracking-wide',
  },
})
```

> [!NOTE]
> La fuente ejecutable es [`packages/ui/src/uno.config.ts`](../packages/ui/src/uno.config.ts). Si este documento y el archivo divergen, **el archivo manda** y este documento debe actualizarse.

**Valores espejo fuera de UnoCSS** (deben mantenerse sincronizados manualmente):
- `packages/ui/src/assets/main.css`: `background-color: #FFFDF5` (`paper`) y `color: #3B3B39` (`ink`).
- `local-app/main.go`: `BackgroundColour: application.NewRGB(255, 253, 245)` (`paper`).

---

## 6. Reglas Inquebrantables de UI para Agentes Constructores

1. **Cero Colores Fuera del Tema:** Prohibido `bg-red-500`, `text-gray-600`, `border-slate-300`, `shadow-lg` o hex sueltos. Solo `paper`, `surface`, `ink`, `brand`, `info`, `success`, `warning` (con opacidades) y sombras `soft*`.
2. **Sin Trazos Duros:** Prohibido `border-2 border-ink`, sombras sólidas desplazadas o negro puro. Bordes siempre por opacidad de `ink`.
3. **Tipografía Única (`Satoshi`) con `tabular-nums` en Contadores:** No importes fuentes adicionales. Temporizadores, PIN, racha y XP llevan siempre `tabular-nums`.
4. **Reka UI sin Estilos Intrusivos:** Modales, menús, tabs y tooltips usan primitivas de **Reka UI** estilizadas con nuestros shortcuts (`card-editorial`, `shadow-soft-lg`, `btn-*`).
5. **Prohibido Inventar Propiedades Visuales No Respaldadas por el Backend:** Cada badge, imagen o contador se alimenta exclusivamente de los tipos de `@bindings/palabrario/pkg/models` a través de `IPalabrarioAPI`.
