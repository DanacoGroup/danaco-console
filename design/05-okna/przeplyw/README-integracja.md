# Animacje 3D kart środowisk pracy — Danaco Console

Cztery komponenty animowane dla kart **TalkIn**, **WorkSpace**, **CodeStudio**, **MultitaskingAI**.

## Pliki

| Plik | Zawartość |
|---|---|
| `danaco-anim-3d.css` | arkusz z wszystkimi czterema animacjami |
| `danaco-anim-3d.html` | fragmenty HTML do wklejenia (po jednym na kartę) |
| `demo.html` | samodzielny podgląd: widok ciemny i jasny, rozmiar karty i powiększony |
| `build_demo.py` | regeneruje `demo.html` po zmianach w CSS/HTML |

## Parametry techniczne

- pętla **3 s**, bezszwowa (wszystkie okresy składowe dzielą 3 s: 3 s / 1,5 s / 1 s)
- **brak tła** — kontener jest w pełni przezroczysty; kolory konturów pochodzą z `currentColor`, więc komponent adaptuje się do widoku ciemnego i jasnego bez przełączników
- czysty CSS 3D (`transform-style: preserve-3d`, `perspective`) — **bez JS i bez zależności**, kompozycja na GPU
- 6–13 elementów DOM na animację, brak przemalowań layoutu
- `prefers-reduced-motion: reduce` → animacja zatrzymana na czytelnej pozie statycznej
- `aria-hidden="true"` / `role="presentation"` — element dekoracyjny, pomijany przez czytniki ekranu

## Integracja

1. Dołącz arkusz (import w globalnym CSS lub `<link>`):

```css
@import "danaco-anim-3d.css";
```

2. Wstaw fragment z `danaco-anim-3d.html` w miejsce statycznej ikony w nagłówku karty środowiska:

| Karta | Klasa | Motyw |
|---|---|---|
| TalkIn | `d3-anim d3-talkin` | orbita dymków rozmowy, wskaźnik pisania |
| WorkSpace | `d3-anim d3-workspace` | izometryczne kafelki zadań, fala unoszenia, znacznik ukończenia |
| CodeStudio | `d3-anim d3-codestudio` | obracający się sześcian kodu, kursor terminala |
| MultitaskingAI | `d3-anim d3-multitasking` | rdzeń orkiestracji, dwie przeciwbieżne orbity węzłów |

## Sterowanie

Wszystko przez zmienne CSS na elemencie animacji lub na dowolnym rodzicu:

```html
<div class="d3-anim d3-talkin" style="--d3-size: 120px"> … </div>
```

| Zmienna | Domyślnie | Znaczenie |
|---|---|---|
| `--d3-size` | `88px` | rozmiar całego komponentu (wszystko skaluje się proporcjonalnie) |
| `--d3-dur` | `3s` | długość pętli |
| `--d3-accent` | per komponent | kolor główny |
| `--d3-accent-2` | per komponent | kolor uzupełniający |

Kolory domyślne: TalkIn `#5b8cff` / `#7fd2ff`, WorkSpace `#17a97d` / `#4fd1a5`,
CodeStudio `#e8883c` / `#ffc078`, MultitaskingAI `#7c6cf0` / `#4bc3ff`.

Podmiana na tokeny projektu, np.:

```css
.card--talkin .d3-anim { --d3-accent: var(--color-accent); --d3-accent-2: var(--color-accent-soft); }
```

## Uwagi wdrożeniowe

- Kontener ma `pointer-events: none`, więc nie przechwytuje kliknięć w kartę.
- Animacja działa również pod `backdrop-filter` karty; w razie artefaktów na WebView usuń
  `backdrop-filter: blur(2px)` z `.d3-talkin .d3-bubble`.
- `color-mix()` ma fallback w bloku `@supports not (...)` dla starszych WebView.
- Aby zatrzymać animacje w testach wizualnych: `document.getAnimations().forEach(a => a.pause())`.
