import './ladowanie.css';

import { elementGodla } from '../ikony/ikony';

/**
 * Scena wejścia — to, co Operator widzi między przejściem bramki a Centrum
 * dowodzenia.
 *
 * ── czemu scena w ogóle stoi ───────────────────────────────────────────────
 * Nie po to, żeby czekanie wyglądało ładniej: żeby czekanie było widoczne.
 * Punkt wejścia składa aplikację i łączy się z rdzeniem pod przesłoną bramki
 * (`main.ts`), więc po wejściu widać stronę zbudowaną, ale jeszcze pustą —
 * pierwszy odczyt (`home.enter`) idzie do rdzenia dopiero wtedy. Bez sceny
 * Operator patrzy przez ten czas na wykaz środowisk, którego nie ma, i nie wie,
 * czy platforma czeka na rdzeń, czy niczego nie znalazła.
 *
 * Scena nie jest więc ładowaniem sztucznym. Stoi nad czasem, który i tak mija,
 * i schodzi w chwili, gdy strona główna ma czym stanąć.
 *
 * ── czemu bryła jest z arkusza, a nie z płótna ─────────────────────────────
 * Sześcian obraca się przestrzenią CSS (`transform-style: preserve-3d`).
 * WebGL dałby to samo wrażenie kosztem zależności, warstwy sterowników
 * i ekranu, który u kogoś nie wstanie — a tu nie ma czego renderować poza
 * jedną bryłą ze znakiem marki na ścianach.
 *
 * ── granice ────────────────────────────────────────────────────────────────
 * Scena nie zna kanału, nie wysyła komend i nie pyta rdzenia o nic. Dostaje
 * obietnicę gotowości i zdanie do pokazania; kiedy gotowość nadejdzie —
 * rozstrzyga ten, kto ją podał. Kres czekania jest twardy: gotowość, która nie
 * nadchodzi, nie ma prawa zamienić sceny w zasłonę nad produktem.
 */

/** Najkrótszy czas na scenie, w milisekundach. */
const PROG_MIGNIECIA = 320;

/**
 * Kres czekania na gotowość, w milisekundach.
 *
 * Po nim scena schodzi mimo wszystko. Strona główna znosi odmowę pierwszego
 * odczytu bez gaszenia ekranu (`widok-strony-glownej.ts`), więc zejście sceny
 * nad stroną, która nic nie dostała, jest gorsze o jedno zdanie w wykazie —
 * a scena, która nie schodzi, jest zamknięciem produktu.
 */
const KRES_CZEKANIA = 8000;

export interface OpisSceny {
  /**
   * Gotowość, na którą scena czeka. Obietnica odrzucona znaczy to samo co
   * spełniona: scena schodzi, bo o tym, czy strona ma czym stanąć, rozstrzyga
   * strona, nie ekran nad nią.
   */
  gotowosc: Promise<unknown>;
  /** Zdanie pod znakiem — czynność, która trwa. */
  krok?: string;
}

export interface ScenaWejscia {
  element: HTMLElement;
  /** Stawia scenę nad dokumentem. */
  uruchom(): void;
  /** Zdejmuje scenę bez czekania — dla ścieżek, które sceny nie potrzebują. */
  zdejmij(): void;
}

/** Buduje scenę wejścia. Nie montuje jej — montaż należy do `uruchom`. */
export function utworzScenaWejscia(opis: OpisSceny): ScenaWejscia {
  const bryla = document.createElement('div');
  bryla.className = 'la-bryla';
  // Znak marki na każdej ścianie: bryła obraca się w kółko, więc ściana bez
  // znaku byłaby co drugi obrót pustym prostokątem.
  for (const strona of ['przod', 'tyl', 'prawo', 'lewo', 'gora', 'dol']) {
    const sciana = document.createElement('div');
    sciana.className = `la-sciana la-sciana--${strona}`;
    // Godło bez etykiety: zdanie o marce pada raz, w napisie pod bryłą.
    // Sześć etykiet znaczyłoby dla czytnika ekranu sześć znaków marki.
    sciana.append(elementGodla({ rozmiar: 48, podloze: 'jasne' }));
    bryla.append(sciana);
  }

  const pole = document.createElement('div');
  pole.className = 'la-pole';
  pole.append(bryla);

  const nazwa = document.createElement('p');
  nazwa.className = 'la-nazwa';
  nazwa.textContent = 'Danaco Console';

  const krok = document.createElement('p');
  krok.className = 'la-krok';
  krok.textContent = opis.krok ?? 'Przygotowuję Centrum dowodzenia…';

  const napisy = document.createElement('div');
  napisy.className = 'la-napisy';
  napisy.append(nazwa, krok);

  const element = document.createElement('div');
  element.className = 'la-scena';
  // Scena mówi o sobie czytnikowi ekranu jednym zdaniem i mówi je na żywo:
  // `status` czyta zmianę treści bez przerywania Operatorowi.
  element.setAttribute('role', 'status');
  element.setAttribute('aria-live', 'polite');
  element.append(pole, napisy);

  let zdjeta = false;

  /** Zdejmuje scenę: wygaszenie, potem usunięcie z dokumentu. */
  function zdejmij(): void {
    if (zdjeta) return;
    zdjeta = true;
    element.classList.add('la-scena--schodzi');
    // Usunięcie idzie po wygaszeniu, a nie zamiast niego. `transitionend` sam
    // nie wystarcza: przy wyłączonym ruchu (`prefers-reduced-motion`) przejście
    // nie zachodzi wcale i zdarzenie nie pada, więc scena zostawałaby na wieki.
    const usun = (): void => element.remove();
    element.addEventListener('transitionend', usun, { once: true });
    window.setTimeout(usun, 400);
  }

  return {
    element,

    uruchom() {
      if (zdjeta || element.isConnected) return;
      document.body.append(element);

      const poczatek = Date.now();
      const zejdz = (): void => {
        // Próg mignięcia: gotowość, która przychodzi natychmiast (rdzeń na tej
        // samej maszynie), dałaby błysk sceny przez dwie klatki. Błysk czyta się
        // jak usterka obrazu, więc scena albo stoi chwilę, albo nie staje wcale.
        const zostalo = PROG_MIGNIECIA - (Date.now() - poczatek);
        if (zostalo > 0) window.setTimeout(zdejmij, zostalo);
        else zdejmij();
      };

      void Promise.race([
        Promise.resolve(opis.gotowosc).catch(() => undefined),
        new Promise((spelnij) => window.setTimeout(spelnij, KRES_CZEKANIA)),
      ]).then(zejdz);
    },

    zdejmij,
  };
}
