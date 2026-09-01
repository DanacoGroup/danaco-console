// Szkielet stron kanału wydań (`pobierz.danaco-group.pl`) — nagłówek, podnawigacja,
// stopka i reguły układu.
//
// GRANICA DWÓCH KANONÓW — TU JEST NAJWAŻNIEJSZE ZDANIE TEGO PLIKU.
//
// Kanon `design/KANON.md` (monochromatyczna precyzja, żetony `--dn-*`, Space
// Grotesk i IBM Plex) rządzi PRODUKTEM i tym, co widać we WNĘTRZU ramek
// z prototypami. NIE rządzi tą witryną. Ta witryna jest stroną na
// `danaco-group.pl` — subdomena to nadal tamta witryna — więc idzie w kanonie
// grupy: jej arkusz `assets/css/styles.css`, jej nagłówek i stopka wzięte
// dosłownie, złote akcenty, kursywa szeryfowa, jej siatka.
//
// Dlaczego to jest twarda granica, a nie preferencja: do tej strony prowadzi
// pozycja POBIERZ z górnego paska witryny firmowej. Człowiek trafia tu jednym
// kliknięciem z `danaco-group.pl`, a zmiana języka wizualnego w pół kroku czyta
// się jako pomyłka albo obca strona — w najlepszym razie jako niedokończona
// robota. Prototypy w ramkach zachowują wygląd Console i to JEST poprawne: one
// pokazują produkt, nie stronę, i mają własne `fundament.css` oraz `prototyp.css`
// zamknięte wewnątrz ramki, więc dwa kanony nie mieszają się w jednym drzewie
// stylów. To cała ich styczność.
//
// PRAKTYCZNIE: poza wnętrzem ramki nie ma tu ani jednego żetonu `--dn-*`, ani
// jednego kroju z `design/zasoby/fonty/`, ani jednej barwy z kanonu Console.
// Wszystko, co poniżej, jest albo zmiennią arkusza grupy (`--space-*`,
// `--blue-950`, `--brand-yellow`, `--text-muted`, `--border-default`,
// `--radius-*`, `--ff-*`), albo regułą układu bez barwy własnej.
//
// KANONU GRUPY NIE MA W POSTACI OPRACOWANIA. Nie leży ani w `design/`, ani
// w `docs/` — jedynym źródłem jest żywa witryna. Wszystko, co tu o nim przyjęto,
// wyprowadzono z jej arkusza i markupu (`portfolio/wytnij-szkielet.mjs` wycina
// nagłówek i stopkę z pobranej strony, nie z przepisania). Gdyby powstała księga
// marki grupy, ona będzie źródłem — nie ten odczyt.
import { readFile } from 'node:fs/promises';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const KATALOG = dirname(fileURLToPath(import.meta.url));

// Metryka z `design/KANON.md`, rozdz. 0 oraz `docs/README.md` — nie z domysłu.
const PRODUCENT = 'Danaco Holding Group Sp. z o.o.';
const TWORCA = 'Dariusz Naharnowicz';
const KONTAKT = 'support@danaco-group.pl';

/** Zamienia znaki o znaczeniu w HTML na encje. Treść stron piszemy po polsku
 *  i z cudzysłowami drukarskimi, więc ucieczka musi być, a nie „chyba nie trzeba". */
export function tekst(wartosc) {
  return String(wartosc)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;');
}

// Własnych reguł jest tyle, ile trzeba na rzeczy, których arkusz grupy nie zna,
// bo dotąd ich na tej witrynie nie było: podnawigacja produktu, karta pakietu
// z sumą kontrolną, scena z prototypem. Wszystko inne — nagłówki sekcji, karty,
// przyciski, tekst ciągły — bierze jej klasy (`section`, `container`, `area`,
// `tag`, `prose`, `btn`, `kicker`, `section-head`, `secno`, `split`, `tlink`).
//
// Krój danych: rodzina systemowa o stałej szerokości znaku. Sumy SHA-256 i liczby
// bajtów muszą stać w kolumnie, a arkusz grupy nie ma żetonu kroju mono.
// Wzięcie tu IBM Plex Mono z dostawy Console byłoby wniesieniem tamtego kanonu
// na tę stronę — więc stoi rodzina systemowa, nie krój marki produktu.
const STYL = `<style>
.dc-dane { font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }

/* MOTYW CIEMNY — POMIAR, NIE ZAŁOŻENIE.
   Arkusz witryny grupy niesie oba motywy, ale dwie jego reguły globalne wiążą
   barwę na stałe: strong oraz h3,h4,h5 dostają w nim color: var(--ink),
   a zmienna --ink (#18233B) NIE przełącza się z motywem. Na jasnym tle to dobra
   barwa; na ciemnym daje granat na granacie — zmierzone zrzutem: nagłówki kart
   pakietów i pogrubienia w uwadze robiły się nieczytelne.
   Poprawiamy to WYŁĄCZNIE wewnątrz własnych komponentów i WYŁĄCZNIE ich własnymi
   żetonami motywu (--text-strong), nie własną barwą — arkusza grupy nie
   nadpisujemy globalnie, bo to ich witryna i ich reguła. */
.dc-pakiet h3, .dc-uwaga strong, .dc-pakiet strong, .dc-tabela strong { color: var(--text-strong); }
.dc-uwaga, .dc-tabela td { color: var(--text-body); }

/* Podnawigacja produktu — pasek pod nagłówkiem witryny. Menu grupy nie zna
   dziesięciu stron kanału i nie ma powodu znać; bez tego paska byłyby osiągalne
   wyłącznie odnośnikami w treści, czyli po omacku. */
/* NAGŁÓWEK WITRYNY GRUPY JEST POZYCJONOWANY NA STAŁE (position: fixed) — i to
   jest rzecz, której nie widać w markupie, a rozstrzyga o całej stronie. Ich
   podstrony kompensują to wysokim wypełnieniem sekcji .page-head; strony kanału
   takiej sekcji nie mają, więc bez tego odstępu pierwszy ekran wchodziłby POD
   nagłówek i pasek stron produktu byłby niewidoczny. Odstęp bierze ich własną
   zmienną wysokości paska, nie wpisaną liczbę — gdy oni ją zmienią, strona
   pójdzie za nimi. Próg 1100 px to ich próg przełączenia na pasek mobilny (zmierzony w arkuszu). */
.dc-podnaw { margin-top: var(--nav-row-height, 152px); background: var(--bg-subtle); border-bottom: 1px solid var(--border-default); }
@media (max-width: 1100px) { .dc-podnaw { margin-top: var(--nav-row-height-mobile, 88px); } }
.dc-podnaw .container { display: flex; align-items: center; gap: var(--space-md); flex-wrap: wrap; padding-top: 10px; padding-bottom: 10px; }
.dc-podnaw__ttl { font-family: var(--ff-display); font-weight: 700; font-size: .84rem; letter-spacing: .06em; text-transform: uppercase; color: var(--text-strong); }
.dc-podnaw a { color: var(--text-muted); text-decoration: none; font-size: .86rem; padding: 2px 0; }
.dc-podnaw a:hover { color: var(--text-strong); }
/* Strona bieżąca nie jest oznaczona samym kolorem — dochodzi wstęga i grubość. */
.dc-podnaw a[aria-current="page"] { color: var(--text-strong); font-weight: 600; box-shadow: inset 0 -2px 0 var(--brand-yellow); }

/* --- Karta pakietu wydania ------------------------------------------------- */
.dc-pakiety { display: grid; gap: var(--space-md); grid-template-columns: repeat(auto-fit, minmax(19rem, 1fr)); margin: var(--space-lg) 0; }
.dc-pakiet { background: var(--bg-surface); border: 1px solid var(--border-default); border-radius: var(--radius-lg); padding: var(--space-md); display: flex; flex-direction: column; gap: 10px; }
.dc-pakiet--twoj { border-color: var(--brand-yellow); box-shadow: var(--shadow-sm); }
.dc-pakiet h3 { margin: 0; font-family: var(--ff-display); font-size: 1.02rem; }
.dc-pakiet p { margin: 0; font-size: .88rem; color: var(--text-muted); }
.dc-pakiet__plik { color: var(--text-strong) !important; font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: .78rem !important; word-break: break-all; }
.dc-pakiet__suma { margin-top: auto !important; font-size: .72rem !important; }
.dc-pakiet__suma code, .dc-pakiet__brak code { display: inline-block; word-break: break-all; background: none; padding: 0; font-size: inherit; }
.dc-pakiet__znacznik { color: var(--text-strong) !important; font-weight: 600; }
.dc-pakiet__haslo { font-size: .76rem !important; }
/* Zdanie o braku podpisu: ta sama wstęga co uwaga o braku — brak nazwany zdaniem, nie samą barwą. */
.dc-pakiet__podpis { font-size: .76rem !important; color: var(--text-body) !important; border-left: 3px solid var(--brand-yellow); padding-left: 10px; }
.dc-pakiet--przygotowanie { border-style: dashed; background: transparent; }
.dc-rozpoznanie { color: var(--text-muted); font-size: .9rem; margin: 6px 0 0; }

/* --- Uwaga ---------------------------------------------------------------- */
.dc-uwaga { border: 1px solid var(--border-default); border-left: 3px solid var(--blue-950); background: var(--bg-subtle); padding: var(--space-md); border-radius: var(--radius-md); margin: var(--space-lg) 0; }
/* Uwaga o braku ma inny nośnik niż barwa: złota wstęga ORAZ zdanie nazywające
   brak na początku treści. Kanon produktu zabrania stanu samym kolorem i ta
   zasada jest tu słuszna niezależnie od tego, czyj arkusz rysuje ramkę. */
.dc-uwaga--brak { border-left-color: var(--brand-yellow); background: var(--bg-sunken); }
.dc-uwaga p { margin: 0; max-width: 82ch; }
.dc-uwaga p + p { margin-top: 10px; }

/* --- Tabela chronologii --------------------------------------------------- */
.dc-tabela { overflow-x: auto; margin: var(--space-lg) 0; }
.dc-tabela table { border-collapse: collapse; width: 100%; font-size: .82rem; }
.dc-tabela th, .dc-tabela td { text-align: left; padding: 10px 12px; border-bottom: 1px solid var(--border-default); vertical-align: top; }
.dc-tabela th { font-family: var(--ff-display); font-size: .72rem; letter-spacing: .08em; text-transform: uppercase; color: var(--text-muted); border-bottom-color: var(--border-strong); white-space: nowrap; }
.dc-tabela code { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: .72rem; word-break: break-all; }

/* --- Zrzut i scena z prototypem ------------------------------------------- */
.dc-zrzuty { display: grid; gap: var(--space-lg); margin: var(--space-lg) 0; }
.dc-zrzuty--para { grid-template-columns: repeat(auto-fit, minmax(24rem, 1fr)); }
.dc-zrzut { margin: 0; }
.dc-zrzut img { display: block; width: 100%; height: auto; border: 1px solid var(--border-default); border-radius: var(--radius-md); }
.dc-zrzut figcaption { margin-top: 8px; color: var(--text-muted); font-size: .78rem; }
</style>`;

/**
 * Składa jedną stronę: szkielet witryny grupy + treść.
 *
 * Szablon czytany jest z pliku, nie wpisany tutaj — powstaje poleceniem
 * `portfolio/wytnij-szkielet.mjs` z pobranej strony witryny grupy, żeby nagłówek
 * i stopka nie były przepisane wzrokiem i nie zestarzały się po cichu.
 */
let szablon = null;

export async function wczytajSzkielet() {
  if (szablon === null) {
    szablon = await readFile(join(KATALOG, '..', 'portfolio', 'szkielet-kanalu.html'), 'utf8');
  }
  return szablon;
}

export function szkielet(strona, wszystkie, szablonKodu) {
  const trasy = wszystkie
    .map(
      (s) =>
        `    <a href="${s.adres}"${s.adres === strona.adres ? ' aria-current="page"' : ''}>${tekst(s.nazwa)}</a>`,
    )
    .join('\n');

  // Treść stron pisana jest zwykłymi znacznikami (`h1`, `h2`, `p`, tabela), bez
  // powtarzania w każdej z nich obudowy sekcji. Obudowa jest tutaj, jedna: arkusz
  // witryny grupy nadaje szerokość, marginesy i tło dopiero wewnątrz
  // `.section > .container` — treść wstawiona wprost do `<main>` wyszłaby
  // przyklejona do lewej krawędzi okna na całą jego szerokość.
  const obudowa = `  <section class="section section--paper">
    <div class="container">
${strona.tresc}
    </div>
  </section>`;

  const stopkaProduktu = `
  <section class="section section--paper">
    <div class="container">
      <div class="prose" style="max-width:70ch;">
        <p style="font-size:.86rem;color:var(--text-muted);">
        <strong>${tekst(PRODUCENT)}</strong> — producent. Twórca: ${tekst(TWORCA)}.
        Kontakt w sprawie produktu i dostępu do wydań:
        <a href="mailto:${tekst(KONTAKT)}">${tekst(KONTAKT)}</a>.
        Status produktu: deweloperski.</p>
      </div>
    </div>
  </section>`;

  return szablonKodu
    .replaceAll('{{TYTUL}}', `${tekst(strona.tytul)} — Danaco Console | Danaco Group`)
    .replaceAll('{{OPIS}}', tekst(strona.opis))
    .replaceAll('{{ADRES}}', tekst(strona.adres))
    .replace('{{STYL}}', STYL)
    .replace('{{TRASY}}', trasy)
    .replace('{{TRESC}}', obudowa + stopkaProduktu);
}
