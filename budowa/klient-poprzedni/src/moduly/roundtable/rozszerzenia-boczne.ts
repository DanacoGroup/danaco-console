import { przyciskAkcji } from '../../modele/kontrolki-formularza-braki';

/** Rozszerzenia boczne modułu Roundtable, warstwa druga widoczności: kody rozszerzeń, te same, którymi okna podpisują się w atrybucie okna. */
export type KodRozszerzenia =
  | 'argument-map-analysis'
  | 'voting-evaluation-center'
  | 'moderator-panel'
  | 'consensus-panel';

/** Okno warstwy drugiej wraz z nazwą, po której Operator je otwiera, i elementem gotowym do osadzenia w powłoce pasa. */
export interface OpisRozszerzenia {
  kod: KodRozszerzenia;
  /** Nazwa własna okna — dokładnie jak w katalogu okien operacyjnych. */
  nazwa: string;
  element: HTMLElement;
}

export interface PasRozszerzen {
  /** Pas osadzany w obszarze modułu pod warstwą pierwszą. */
  element: HTMLElement;
  /** Przycisk otwierający rozszerzenie; jego stan idzie za stanem pasa, także gdy otwarto skądinąd. */
  przyciskOtwarcia(kod: KodRozszerzenia): HTMLButtonElement;
  /** Otwiera rozszerzenie i przenosi do niego ognisko. */
  otworz(kod: KodRozszerzenia): void;
  /** Zwija rozszerzenie; treść zostaje w dokumencie, nasłuch zostaje żywy. */
  zwin(kod: KodRozszerzenia): void;
}

/** Stan jednego rozszerzenia wraz z jego powłoką, przyciskami otwarcia i nastawą, czy jest teraz otwarte. */
interface WpisRozszerzenia {
  opis: OpisRozszerzenia;
  powloka: HTMLElement;
  przyciski: HTMLButtonElement[];
  otwarte: boolean;
}

/**
 * Licznik pasów w dokumencie — identyfikatory powłok muszą być niepowtarzalne,
 * bo `aria-controls` wiąże przycisk z powłoką po identyfikatorze, a moduł
 * i panel pomocniczy Roundtable potrafią stać w dokumencie jednocześnie.
 */
let licznikPasow = 0;

export function utworzPasRozszerzen(rozszerzenia: readonly OpisRozszerzenia[]): PasRozszerzen {
  const wpisy = new Map<KodRozszerzenia, WpisRozszerzenia>();
  licznikPasow += 1;
  const numerPasa = licznikPasow;

  const element = document.createElement('div');
  element.className = 'dr-modul__rozszerzenia';
  element.setAttribute('aria-label', 'Rozszerzenia boczne modułu Roundtable');

  for (const opis of rozszerzenia) {
    const powloka = document.createElement('div');
    powloka.className = 'dr-rozszerzenie';
    powloka.id = `dr-rozszerzenie-${numerPasa}-${opis.kod}`;
    powloka.dataset['rozszerzenie'] = opis.kod;
    powloka.hidden = true;

    const zwiniecie = przyciskAkcji(`Zwiń rozszerzenie: ${opis.nazwa}`, 'dn-btn dn-btn--zarys dn-btn--sm');
    zwiniecie.addEventListener('click', () => ustaw(opis.kod, false));

    const pasek = document.createElement('div');
    pasek.className = 'dr-rozszerzenie__pasek';
    pasek.append(zwiniecie);

    powloka.append(pasek, opis.element);
    element.append(powloka);
    wpisy.set(opis.kod, { opis, powloka, przyciski: [], otwarte: false });
  }

  // Wpis rozszerzenia albo błąd wprost: kod nieznany pasowi jest usterką złożenia, nie Operatora.
  function wpisAlbo(kod: KodRozszerzenia): WpisRozszerzenia {
    const wpis = wpisy.get(kod);
    if (wpis === undefined) {
      throw new Error(`Pas rozszerzeń modułu Roundtable nie zna okna o kodzie ${kod}.`);
    }
    return wpis;
  }

  /** Jedno miejsce zmiany stanu — powłoka, przyciski i ognisko naraz. */
  function ustaw(kod: KodRozszerzenia, otwarte: boolean): void {
    const wpis = wpisAlbo(kod);
    wpis.otwarte = otwarte;
    wpis.powloka.hidden = !otwarte;
    for (const przycisk of wpis.przyciski) {
      przycisk.setAttribute('aria-expanded', String(otwarte));
      przycisk.dataset['otwarte'] = otwarte ? 'tak' : 'nie';
      przycisk.textContent = etykietaPrzycisku(wpis.opis.nazwa, otwarte);
    }
    if (!otwarte) return;
    wpis.opis.element.focus();
    wpis.opis.element.scrollIntoView({ block: 'nearest' });
  }

  return {
    element,

    przyciskOtwarcia(kod) {
      const wpis = wpisAlbo(kod);
      const przycisk = przyciskAkcji(
        etykietaPrzycisku(wpis.opis.nazwa, wpis.otwarte),
        'dn-btn dn-btn--zarys',
      );
      przycisk.setAttribute('aria-expanded', String(wpis.otwarte));
      przycisk.setAttribute('aria-controls', wpis.powloka.id);
      przycisk.dataset['otwarte'] = wpis.otwarte ? 'tak' : 'nie';
      przycisk.dataset['rozszerzenie'] = kod;
      przycisk.addEventListener('click', () => ustaw(kod, !wpisAlbo(kod).otwarte));
      wpis.przyciski.push(przycisk);
      return przycisk;
    },

    otworz: (kod) => ustaw(kod, true),
    zwin: (kod) => ustaw(kod, false),
  };
}

/** Etykieta przycisku niesie stan słowem, nie samą strzałką ani samą barwą, więc czytelnik ekranu ją usłyszy. */
function etykietaPrzycisku(nazwa: string, otwarte: boolean): string {
  return otwarte ? `Zwiń: ${nazwa}` : `Otwórz: ${nazwa}`;
}
