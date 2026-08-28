import type { Kwit } from './kwit-decyzji';
import type { PozycjaDecyzji } from './pozycje-decyzji';
import {
  ETYKIETY_NASTAW,
  ETYKIETY_ZASIEGOW,
  type NastawaKoordynatora,
  type WywolaniaInterwencji,
} from './wywolania-interwencji';

/**
 * Arkusz dróg — drugie dotknięcie: karta pozycji otwiera arkusz, przycisk
 * w nim wykonuje drogę. Przycisk powstaje tylko tam, gdzie droga jest
 * przejezdna. Przycisk przejęcia zostaje czynny nawet z pustym polem.
 */
export interface ArkuszDrog {
  element: HTMLElement;
  pokaz(pozycja: PozycjaDecyzji): void;
  ukryj(): void;
  /** Czy arkusz jest przywołany — do sprawdzianów i do obsługi klawisza wyjścia. */
  otwarty(): boolean;
}

export interface OpisArkusza {
  wywolania: WywolaniaInterwencji;
  /** Kwit każdej drogi idzie tędy do paska kwitu i do pamięci trwałej. */
  poKwicie: (kwit: Kwit) => void;
}

/** Trzy gotowe nastawy koordynatora dostępne w arkuszu dróg, uporządkowane od tej najczęściej wybieranej. */
const NASTAWY: readonly NastawaKoordynatora[] = [
  'pytaj-o-kazdy-krok',
  'stoj-przy-braku-postepu',
  'obniz-naklad',
];

export function utworzArkuszDrog(opis: OpisArkusza): ArkuszDrog {
  const naglowek = document.createElement('h4');
  naglowek.className = 'mb-arkusz__naglowek';

  const kontekst = document.createElement('ul');
  kontekst.className = 'mb-arkusz__kontekst';

  const drogi = document.createElement('div');
  drogi.className = 'mb-arkusz__drogi';

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn dn-btn--zarys mb-cel';
  zamknij.textContent = 'Wróć do wykazu';
  zamknij.addEventListener('click', () => ukryj());

  const element = document.createElement('section');
  element.className = 'mb-arkusz';
  element.hidden = true;
  element.setAttribute('aria-label', 'Drogi interwencji dla wybranej pozycji');
  element.append(naglowek, kontekst, drogi, zamknij);

  function ukryj(): void {
    element.hidden = true;
    drogi.replaceChildren();
  }

  /** Przycisk drogi — jedno dotknięcie, jedno wywołanie, jeden kwit. */
  function przyciskDrogi(etykieta: string, odmiana: string, wykonaj: () => Promise<Kwit>): HTMLButtonElement {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = `dn-btn ${odmiana} mb-cel mb-arkusz__droga`;
    przycisk.textContent = etykieta;
    przycisk.addEventListener('click', () => {
      // Znacznik pracy zamiast zablokowania: kontrolka zostaje czynna.
      przycisk.setAttribute('aria-busy', 'true');
      void wykonaj().then((kwit) => {
        przycisk.removeAttribute('aria-busy');
        opis.poKwicie(kwit);
        ukryj();
      });
    });
    return przycisk;
  }

  /** Zdanie o drodze nieprzejezdnej — na miejscu przycisku, którego nie ma. */
  function zdanieBraku(tresc: string): HTMLElement {
    const akapit = document.createElement('p');
    akapit.className = 'mb-arkusz__brak';
    akapit.textContent = tresc;
    return akapit;
  }

  function sekcja(tytul: string): HTMLElement {
    const naglowekSekcji = document.createElement('h5');
    naglowekSekcji.className = 'mb-arkusz__sekcja';
    naglowekSekcji.textContent = tytul;
    return naglowekSekcji;
  }

  return {
    element,

    pokaz(pozycja) {
      naglowek.textContent = pozycja.naglowek;
      kontekst.replaceChildren(
        ...pozycja.kontekst.map((zdanie) => {
          const wiersz = document.createElement('li');
          wiersz.textContent = zdanie;
          return wiersz;
        }),
      );

      const dostepne = opis.wywolania.drogiDostepne(pozycja);
      const czesci: HTMLElement[] = [];

      czesci.push(sekcja('1. Zatwierdź krok'));
      if (dostepne.zatwierdz.dostepna) {
        czesci.push(
          przyciskDrogi('Puść krok dalej', 'dn-btn--sygnal', () =>
            opis.wywolania.zatwierdzKrok(pozycja, 'pusc-dalej'),
          ),
          przyciskDrogi('Powtórz krok', 'dn-btn--zarys', () =>
            opis.wywolania.zatwierdzKrok(pozycja, 'powtorz'),
          ),
        );
      } else {
        czesci.push(zdanieBraku(dostepne.zatwierdz.powod ?? ''));
      }

      czesci.push(sekcja('2. Wstrzymaj'));
      if (dostepne.wstrzymaj.dostepna) {
        for (const zasieg of dostepne.wstrzymaj.zasiegi) {
          czesci.push(
            przyciskDrogi(ETYKIETY_ZASIEGOW[zasieg], 'dn-btn--atrament', () =>
              opis.wywolania.wstrzymaj(pozycja, zasieg),
            ),
          );
        }
      } else {
        czesci.push(zdanieBraku(dostepne.wstrzymaj.powod ?? ''));
      }

      czesci.push(sekcja('3. Nastaw koordynatora'));
      if (dostepne.nastaw.dostepna) {
        for (const nastawa of NASTAWY) {
          czesci.push(
            przyciskDrogi(ETYKIETY_NASTAW[nastawa], 'dn-btn--zarys', () =>
              opis.wywolania.nastawKoordynatora(pozycja, nastawa),
            ),
          );
        }
      } else {
        czesci.push(zdanieBraku(dostepne.nastaw.powod ?? ''));
      }

      czesci.push(sekcja('4. Przejmij bezpośrednie sterowanie'));
      if (dostepne.przejmij.dostepna) {
        const wyjasnienie = document.createElement('p');
        wyjasnienie.className = 'mb-arkusz__wyjasnienie';
        wyjasnienie.textContent =
          'Zatrzymuję turę, przestawiam okno na pytanie o każdy krok i wpuszczam twoje ' +
          'polecenie. Dedykowanej komendy przejęcia kontrakt dziś nie ma — tak wygląda ' +
          'przejęcie złożone z trzech wywołań, które rdzeń naprawdę wykonuje.';

        const polecenie = document.createElement('textarea');
        polecenie.className = 'dn-pole-kontrolka mb-arkusz__polecenie';
        polecenie.rows = 3;
        polecenie.placeholder = 'Twoje polecenie dla okna…';
        polecenie.setAttribute('aria-label', 'Polecenie Operatora wpuszczane do okna');

        czesci.push(
          wyjasnienie,
          polecenie,
          przyciskDrogi('Przejmij i wyślij polecenie', 'dn-btn--niebezpieczny', () =>
            opis.wywolania.przejmij(pozycja, polecenie.value),
          ),
        );
      } else {
        czesci.push(zdanieBraku(dostepne.przejmij.powod ?? ''));
      }

      drogi.replaceChildren(...czesci);
      element.hidden = false;
    },

    ukryj,
    otwarty: () => !element.hidden,
  };
}
