import {
  DiffHunkKind,
  type StudioDiffHunk,
  type StudioSemanticMatch,
} from '../../../../shared/contract';
import { opiszLiczniki, policzTresc } from './liczniki-dokumentu';
import {
  ocenaCzytelnosci,
  opiszOcene,
  policzCzytelnosc,
  pomiaryKorekty,
  uscisleniaTresci,
} from './ocena-redaktora';
import { opiszStatystyke, policzRoznice } from './statystyka-roznicy';

/**
 * Panel Redaktora — ocena dokumentu, korekty, uściślenia, podobieństwa
 * i statystyki, po prawej stronie treści.
 *
 * ── Zasada panelu: każda liczba policzona ───────────────────────────────────
 * Ocena wychodzi z miary czytelności liczonej z tekstu (`ocena-redaktora.ts`),
 * a nie z wrażenia. Korekty liczą wzorce, które da się rozpoznać w zapisie;
 * pisownia i gramatyka stoją w wykazie z wartością NIEPODANĄ i z powodem, bo
 * słownika języka i analizy składniowej rdzeń nie ma, a program spoza instalki
 * jest zakazany. Podobieństwa liczy rdzeń — `studio.diff.source` zestawia
 * dokument z materiałem wejściowym, a `studio.search.semantic` szuka fragmentów
 * bliskich znaczeniowo. Statystyki bierze `liczniki-dokumentu.ts`.
 *
 * Panel nie woła rdzenia sam: okno podaje mu wyniki, które ma. Dzięki temu
 * pomiar zlecony i pomiar policzony na miejscu nie mieszają się w jednym
 * miejscu kodu.
 */

/** Czynności panelu zlecane oknu. */
export interface CzynnosciRedaktora {
  /** Uruchamia uściślenie operacją kontekstową. */
  naUscislenie(idAkcji: string): void;
  /** Zleca pomiar podobieństw — zestawienie z materiałem wejściowym. */
  naPodobienstwa(wskazanie: string): void;
  /** Zleca wyszukiwanie znaczeniowe wskazanym zapytaniem. */
  naWyszukanie(zapytanie: string): void;
}

/** Panel wraz z jego odświeżeniem. */
export interface PanelRedaktora {
  element: HTMLElement;
  /** Przelicza pomiary z treści bieżącej. */
  odswiez(tresc: string): void;
  /** Przyjmuje wynik zestawienia ze źródłem oddany przez rdzeń. */
  ustawPodobienstwa(fragmenty: readonly StudioDiffHunk[], odczytano: boolean): void;
  /** Przyjmuje wynik wyszukiwania znaczeniowego. */
  ustawTrafienia(trafienia: readonly StudioSemanticMatch[], droga: string): void;
  /** Otwiera albo zamyka panel. */
  przestawWidocznosc(): void;
  widoczny(): boolean;
}

export function utworzPanelRedaktora(czynnosci: CzynnosciRedaktora): PanelRedaktora {
  let otwarty = true;

  const ocena = document.createElement('p');
  ocena.className = 'ms-redaktor__ocena';

  const podstawaOceny = document.createElement('p');
  podstawaOceny.className = 'dn-pole-opis';

  const korekty = document.createElement('ul');
  korekty.className = 'ms-redaktor__wykaz';

  const uscislenia = document.createElement('ul');
  uscislenia.className = 'ms-redaktor__wykaz';

  const podobienstwa = document.createElement('div');
  podobienstwa.className = 'ms-redaktor__podobienstwa';

  const statystyki = document.createElement('p');
  statystyki.className = 'dn-pole-opis ms-redaktor__statystyki';

  const wskazanieZrodla = document.createElement('input');
  wskazanieZrodla.type = 'text';
  wskazanieZrodla.className = 'dn-pole-kontrolka';
  wskazanieZrodla.placeholder = 'identyfikator pliku Library — materiał wejściowy';
  wskazanieZrodla.setAttribute('aria-label', 'Materiał wejściowy do zestawienia');

  const zlecPodobienstwa = document.createElement('button');
  zlecPodobienstwa.type = 'button';
  zlecPodobienstwa.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  zlecPodobienstwa.textContent = 'Zestaw ze źródłem';
  zlecPodobienstwa.dataset['czynnosc'] = 'podobienstwa';
  zlecPodobienstwa.addEventListener('click', () =>
    czynnosci.naPodobienstwa(wskazanieZrodla.value.trim()),
  );

  const zapytanie = document.createElement('input');
  zapytanie.type = 'text';
  zapytanie.className = 'dn-pole-kontrolka';
  zapytanie.placeholder = 'czego szukać znaczeniowo w dokumencie';
  zapytanie.setAttribute('aria-label', 'Zapytanie wyszukiwania znaczeniowego');

  const zlecWyszukanie = document.createElement('button');
  zlecWyszukanie.type = 'button';
  zlecWyszukanie.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  zlecWyszukanie.textContent = 'Szukaj znaczeniowo';
  zlecWyszukanie.dataset['czynnosc'] = 'znaczeniowo';
  zlecWyszukanie.addEventListener('click', () => czynnosci.naWyszukanie(zapytanie.value.trim()));

  const element = document.createElement('aside');
  element.className = 'ms-redaktor';
  element.setAttribute('aria-label', 'Panel Redaktora — pomiary dokumentu');
  element.append(
    czesc('Ocena dokumentu', [ocena, podstawaOceny]),
    czesc('Korekty', [korekty]),
    czesc('Uściślenia', [uscislenia]),
    czesc('Podobieństwa', [wskazanieZrodla, zlecPodobienstwa, zapytanie, zlecWyszukanie, podobienstwa]),
    czesc('Statystyki dokumentu', [statystyki]),
  );

  function odswiez(tresc: string): void {
    const liczniki = policzTresc(tresc);
    const czytelnosc = policzCzytelnosc(tresc);

    ocena.textContent = `${ocenaCzytelnosci(czytelnosc)} / 100`;
    ocena.dataset['ocena'] = String(ocenaCzytelnosci(czytelnosc));
    podstawaOceny.textContent = opiszOcene(czytelnosc);

    korekty.replaceChildren(
      ...pomiaryKorekty(tresc).map((pomiar) => {
        const nazwa = document.createElement('strong');
        nazwa.textContent = pomiar.nazwa;

        const wartosc = document.createElement('span');
        wartosc.className = 'dn-plakietka ms-redaktor__licznik';
        wartosc.textContent =
          pomiar.wartosc === null ? 'bez pomiaru' : `${pomiar.wartosc} ${pomiar.miano}`;

        const podstawa = document.createElement('span');
        podstawa.className = 'dn-pole-opis';
        podstawa.textContent = pomiar.podstawa;

        const pozycja = document.createElement('li');
        pozycja.dataset['pomiar'] = pomiar.kod;
        pozycja.dataset['mierzone'] = pomiar.wartosc === null ? 'nie' : 'tak';
        pozycja.append(nazwa, wartosc, podstawa);
        return pozycja;
      }),
    );

    uscislenia.replaceChildren(
      ...uscisleniaTresci(tresc, liczniki).map((uscislenie) => {
        const odhaczenie = document.createElement('input');
        odhaczenie.type = 'checkbox';
        odhaczenie.className = 'dn-przelacznik';
        odhaczenie.dataset['uscislenie'] = uscislenie.kod;
        odhaczenie.title =
          'Odhaczenie jest znacznikiem Operatora — „to miejsce mam już przejrzane". ' +
          'Kontrakt nie ma pola na taki znacznik, więc żyje on przez tę sesję okna.';

        const nazwa = document.createElement('strong');
        nazwa.textContent = uscislenie.nazwa;

        const wskazanie = document.createElement('span');
        wskazanie.className = 'dn-pole-opis';
        wskazanie.textContent = uscislenie.wskazanie;

        const uruchom = document.createElement('button');
        uruchom.type = 'button';
        uruchom.className = 'dn-btn dn-btn--sm dn-btn--duch';
        uruchom.textContent = 'Zleć modelowi';
        uruchom.addEventListener('click', () => czynnosci.naUscislenie(uscislenie.idAkcji));

        const pozycja = document.createElement('li');
        pozycja.dataset['uscislenie'] = uscislenie.kod;
        pozycja.append(odhaczenie, nazwa, wskazanie, uruchom);
        return pozycja;
      }),
    );

    statystyki.textContent = `${opiszLiczniki(liczniki)} · słów długich ${czytelnosc.udzialSlowDlugich} %`;
  }

  return {
    element,
    odswiez,

    ustawPodobienstwa(fragmenty, odczytano) {
      if (!odczytano) {
        podobienstwa.textContent =
          'Rdzeń nie odczytał materiału wejściowego, więc podobieństw nie ma z czym porównać — ' +
          'to odpowiedź, nie brak odpowiedzi (pole sourceResolved komendy studio.diff.source).';
        return;
      }
      if (fragmentyPuste(fragmenty)) {
        podobienstwa.textContent =
          'Rdzeń zestawił dokument z materiałem wejściowym i nie znalazł rozbieżności: treść ' +
          'dokumentu zgadza się ze źródłem w całości. Wysoka zgodność jest tu wskazaniem ' +
          'przejęcia, nie zasługą.';
        return;
      }
      podobienstwa.textContent =
        `${opiszStatystyke(policzRoznice(fragmenty))} Fragmenty kontekstowe to miejsca ZGODNE ` +
        'ze źródłem — one wskazują przejęcie; fragmenty zmienione to praca własna.';
    },

    ustawTrafienia(trafienia, droga) {
      if (trafienia.length === 0) {
        podobienstwa.textContent =
          `Wyszukiwanie znaczeniowe (droga: ${droga}) nie oddało ani jednego fragmentu.`;
        return;
      }
      const wykaz = document.createElement('ul');
      wykaz.className = 'ms-redaktor__wykaz';
      for (const trafienie of trafienia) {
        const pozycja = document.createElement('li');
        const miara = document.createElement('span');
        miara.className = 'dn-plakietka ms-redaktor__licznik';
        miara.textContent = `bliskość ${Math.round(trafienie.score * 100) / 100}`;
        const tekst = document.createElement('span');
        tekst.textContent = trafienie.text;
        pozycja.append(miara, tekst);
        wykaz.append(pozycja);
      }
      podobienstwa.replaceChildren(
        zdanie(`Droga pomiaru: ${droga}. Fragmentów: ${trafienia.length}.`),
        wykaz,
      );
    },

    przestawWidocznosc() {
      otwarty = !otwarty;
      element.hidden = !otwarty;
    },

    widoczny: () => otwarty,
  };
}

/** Czy zestawienie nie ma ani jednego fragmentu zmiany. */
function fragmentyPuste(fragmenty: readonly StudioDiffHunk[]): boolean {
  return fragmenty.every((fragment) => fragment.kind === DiffHunkKind.Context);
}

/** Sekcja panelu wraz z jej tytułem. */
function czesc(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-redaktor__tytul';
  naglowek.textContent = tytul;

  const sekcja = document.createElement('section');
  sekcja.className = 'ms-redaktor__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}

/** Zdanie opisu — jedno miejsce składania akapitu pomocniczego. */
function zdanie(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = tresc;
  return element;
}
