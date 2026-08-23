import { opisOdmowy } from '../../komponenty/odmowa';
import { pole, przyciskAkcji } from '../../modele/kontrolki-formularza';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';

/**
 * Zakładka „Skróty i schowek" Command & Tools Hub — dziewięć komend trzech
 * rodzin: historia schowka, słownik skrótów tekstowych i skrót globalny
 * wywoływacza poleceń.
 *
 * ── Czego ta zakładka NIE robi ──────────────────────────────────────────────
 * Nie czyta schowka maszyny Operatora i nie udaje, że umie. Schowek należy do
 * tamtej maszyny, a rdzeń stoi na serwerze. Podział jest jawny i widoczny na
 * ekranie: Operator wkleja skopiowaną treść w pole „Treść do zapamiętania",
 * okno oddaje ją rdzeniowi (`clipboard.push`), a rdzeń daje jej trwałość.
 * Powrót idzie tą samą drogą: kliknięcie wpisu kopiuje go do schowka karty
 * przeglądarką, a nie rdzeniem.
 *
 * ── Skrót globalny ──────────────────────────────────────────────────────────
 * Nastawę trzyma rdzeń, klawisze przechwytuje powłoka programu okiennego.
 * Odpowiedź mówi wprost, czy rejestracji ma kto dokonać, a okno powtarza to
 * zdanie — zamiast obiecywać skrót, który nikogo nie obudzi.
 */
export interface PanelSchowka {
  element: HTMLElement;
  /** Odczyt historii, słownika skrótów i nastawy skrótu globalnego. */
  wczytaj(): Promise<void>;
}

export function utworzPanelSchowka(stan: StanAssistant): PanelSchowka {
  const okno: StanOkna = utworzStanOkna();

  // ── Historia schowka ──────────────────────────────────────────────────────
  const trescSchowka = pole('Treść do zapamiętania', 'wklej skopiowaną treść');
  const dopisz = przyciskAkcji('Zapamiętaj w rdzeniu', 'dn-btn dn-btn--sm dn-btn--zarys');
  dopisz.addEventListener('click', () => void dopiszTresc());

  const szukaj = pole('Szukaj w historii', 'fragment treści');
  szukaj.addEventListener('change', () => void wczytajHistorie());

  const wyczysc = przyciskAkcji('Wyczyść historię', 'dn-btn dn-btn--sm dn-btn--duch');
  wyczysc.addEventListener('click', () => void wyczyscHistorie());

  const historia = document.createElement('ul');
  historia.className = 'ma-wykaz';
  historia.dataset['wykaz'] = 'schowek';

  // ── Słownik skrótów tekstowych ────────────────────────────────────────────
  const skrot = pole('Skrót', 'np. ;odmowa');
  const rozwiniecie = pole('Rozwinięcie skrótu', 'treść, w którą skrót się rozwija');
  const zapiszSkrot = przyciskAkcji('Zapisz skrót', 'dn-btn dn-btn--sm dn-btn--zarys');
  zapiszSkrot.addEventListener('click', () => void zapiszSkrotTekstowy());

  const skroty = document.createElement('ul');
  skroty.className = 'ma-wykaz';
  skroty.dataset['wykaz'] = 'skroty';

  // ── Skrót globalny wywoływacza ────────────────────────────────────────────
  const hotkey = pole('Skrót globalny wywoływacza', 'np. Ctrl+Shift+Space');
  const zapiszHotkey = przyciskAkcji('Zapisz skrót globalny', 'dn-btn dn-btn--sm dn-btn--zarys');
  zapiszHotkey.addEventListener('click', () => void zapiszWywolywacz());
  const stanHotkey = document.createElement('p');
  stanHotkey.className = 'dn-pole-opis';

  okno.tresc.append(
    sekcja('Historia schowka', [trescSchowka, dopisz, szukaj, wyczysc, historia]),
    sekcja('Słownik skrótów tekstowych', [skrot, rozwiniecie, zapiszSkrot, skroty]),
    sekcja('Wywoływacz poleceń', [hotkey, zapiszHotkey, stanHotkey]),
  );

  const element = document.createElement('div');
  element.className = 'ma-panel';
  element.dataset['panel'] = 'schowek';
  element.append(okno.element);

  async function wczytaj(): Promise<void> {
    okno.ladowanie('Odczyt historii schowka, słownika skrótów i skrótu globalnego…');
    await Promise.all([wczytajHistorie(), wczytajSkroty(), wczytajWywolywacz()]);
  }

  async function wczytajHistorie(): Promise<void> {
    const wynik = await stan.schowek.historia(szukaj.value, false);
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt historii schowka', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    historia.replaceChildren();
    for (const wpis of wynik.wynik.entries) {
      const wiersz = document.createElement('li');
      wiersz.className = 'ma-wykaz__wiersz';
      wiersz.dataset['wpis'] = wpis.id;

      const tresc = document.createElement('span');
      tresc.textContent = wpis.preview ?? wpis.content;

      // Wklejenie należy do karty, nie do rdzenia: to schowek maszyny
      // Operatora, a rdzeń go nie widzi.
      const skopiuj = przyciskAkcji('Kopiuj', 'dn-btn dn-btn--xs dn-btn--duch');
      skopiuj.addEventListener('click', () => {
        void navigator.clipboard?.writeText(wpis.content);
      });

      const przypnij = przyciskAkcji(
        wpis.pinned ? 'Odepnij' : 'Przypnij',
        'dn-btn dn-btn--xs dn-btn--duch',
      );
      przypnij.addEventListener('click', () => void przypnijWpis(wpis.id, !wpis.pinned));

      const usun = przyciskAkcji('Usuń', 'dn-btn dn-btn--xs dn-btn--duch');
      usun.addEventListener('click', () => void usunWpis(wpis.id));

      wiersz.append(tresc, skopiuj, przypnij, usun);
      historia.append(wiersz);
    }
    if (wynik.wynik.entries.length === 0) {
      okno.puste(
        'Historia schowka jest pusta. Rdzeń nie czyta schowka Twojej maszyny — ' +
          'wklej treść w pole powyżej, a rdzeń da jej trwałość.',
      );
      return;
    }
    okno.gotowe();
  }

  async function dopiszTresc(): Promise<void> {
    if (trescSchowka.value.trim() === '') {
      okno.blad('Pole treści jest puste — nie ma czego zapamiętać.');
      return;
    }
    const wynik = await stan.schowek.dopisz({ tresc: trescSchowka.value });
    if (!wynik.udany) {
      okno.blad(opisOdmowy('Zapis wpisu schowka', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    trescSchowka.value = '';
    await wczytajHistorie();
  }

  async function przypnijWpis(id: string, przypiety: boolean): Promise<void> {
    const wynik = await stan.schowek.przypnij(id, przypiety);
    if (!wynik.udany) {
      okno.blad(opisOdmowy('Przypięcie wpisu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    await wczytajHistorie();
  }

  async function usunWpis(id: string): Promise<void> {
    const wynik = await stan.schowek.usun(id, true);
    if (!wynik.udany) {
      okno.blad(opisOdmowy('Usunięcie wpisu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    await wczytajHistorie();
  }

  async function wyczyscHistorie(): Promise<void> {
    const wynik = await stan.schowek.usun('', false);
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Czyszczenie historii', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    await wczytajHistorie();
  }

  async function wczytajSkroty(): Promise<void> {
    const wynik = await stan.schowek.skroty('');
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt słownika skrótów', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    skroty.replaceChildren();
    for (const pozycja of wynik.wynik.snippets) {
      const wiersz = document.createElement('li');
      wiersz.className = 'ma-wykaz__wiersz';
      wiersz.dataset['skrot'] = pozycja.id;

      const opis = document.createElement('span');
      opis.textContent = `${pozycja.shortcut} → ${pozycja.content}`;

      const usun = przyciskAkcji('Usuń', 'dn-btn dn-btn--xs dn-btn--duch');
      usun.addEventListener('click', () => void usunSkrotZeSlownika(pozycja.id));

      wiersz.append(opis, usun);
      skroty.append(wiersz);
    }
  }

  async function zapiszSkrotTekstowy(): Promise<void> {
    const wynik = await stan.schowek.zapiszSkrot({
      id: '',
      skrot: skrot.value.trim(),
      tresc: rozwiniecie.value,
      opis: '',
      polaSzablonu: [],
      czynny: true,
    });
    if (!wynik.udany) {
      okno.blad(opisOdmowy('Zapis skrótu tekstowego', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    skrot.value = '';
    rozwiniecie.value = '';
    await wczytajSkroty();
  }

  async function usunSkrotZeSlownika(id: string): Promise<void> {
    const wynik = await stan.schowek.usunSkrot(id);
    if (!wynik.udany) {
      okno.blad(opisOdmowy('Usunięcie skrótu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    await wczytajSkroty();
  }

  async function wczytajWywolywacz(): Promise<void> {
    const wynik = await stan.schowek.skrotGlobalny();
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt skrótu globalnego', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    hotkey.value = wynik.wynik.hotkey;
    stanHotkey.textContent = zdanieOSkrocie(
      wynik.wynik.supported,
      wynik.wynik.registered,
      wynik.wynik.reason,
    );
  }

  async function zapiszWywolywacz(): Promise<void> {
    const wynik = await stan.schowek.zapiszSkrotGlobalny(hotkey.value.trim());
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Zapis skrótu globalnego', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    // Zapis udaje się także wtedy, gdy rejestracji nie ma kto wykonać — tak
    // stanowi kontrakt. Okno powtarza więc oba rozstrzygnięcia osobno.
    stanHotkey.textContent = zdanieOSkrocie(true, wynik.wynik.registered, wynik.wynik.reason);
  }

  return { element, wczytaj };
}

/** Zdanie o skrócie globalnym: co zapisano i czy ktokolwiek to przechwyci. */
function zdanieOSkrocie(
  wspierany: boolean,
  zarejestrowany: boolean,
  powod: string | undefined,
): string {
  if (zarejestrowany) return 'Skrót jest zapisany i zarejestrowany w powłoce.';
  const koncowka = powod !== undefined && powod !== '' ? ` ${powod}` : '';
  if (!wspierany) return `Skrót jest zapisany, ale nie ma go kto przechwycić.${koncowka}`;
  return `Skrót jest zapisany; rejestracja nie doszła do skutku.${koncowka}`;
}

/** Sekcja panelu: tytuł wraz z kontrolkami jednego obszaru. */
function sekcja(tytul: string, dzieci: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('h4');
  naglowek.className = 'ma-panel__tytul';
  naglowek.textContent = tytul;

  const element = document.createElement('div');
  element.className = 'ma-panel__sekcja';
  element.append(naglowek, ...dzieci);
  return element;
}
