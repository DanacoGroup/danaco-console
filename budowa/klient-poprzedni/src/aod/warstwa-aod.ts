import './aod.css';

import type { Kanal } from '../protokol/kanal';
import { StanAwatara, utworzAwatarAod } from './awatar-aod';
import { utworzDymekSugestii } from './dymek-sugestii';
import type { WpisKolejkiDecyzji } from './kolejka-decyzji';
import { utworzKolumneAod, type KolumnaAod, type OpisKolumnyAod } from './kolumna-aod';
import { WagaUjawnienia } from './rodzaje-sugestii';
import {
  TrybObecnosci,
  utworzStanObecnosci,
  type OpisUjawnienia,
  type StanObecnosci,
} from './tryb-obecnosci';
import { KLASA_POWODU, type MagazynWyciszen } from './wyciszenie-aod';
import { utworzKontekstWyciszenia, type KontekstWyciszenia } from './wyciszenie-kontekst';
import { utworzMenuWyciszenia } from './wyciszenie-menu';

/** Warstwa Always On Display łączy pływający awatar, dymek kontekstowy sugestii i kolumnę boczną powierzchni interakcji w jedną powierzchnię osadzoną poza obszarem modułu, poza jego przełączeniami. */
export interface WarstwaAod {
  /** Element do osadzenia w powłoce, poza obszarem roboczym modułu. */
  element: HTMLElement;
  /** Otwiera powierzchnię interakcji — listwa ustawień i skrót klawiszowy. */
  otworzPowierzchnie(): void;
  /** Stan obecności — do sterowania z okna Ustawień, gdy takie wejście powstanie. */
  obecnosc: StanObecnosci;
  /** Kolumna powierzchni interakcji. */
  kolumna: KolumnaAod;
  rozlacz(): void;
}

/** Stała podaje szerokość obszaru roboczego, poniżej której kolumna nie mieści się obok pracy i powierzchnia interakcji zwija się do dymka kontekstowego ze skróconą treścią. */
const PROG_CZYTELNOSCI = 640;

/** Funkcja przyjmuje magazyn stanu wyciszeń i trybu obecności jako parametr: pominięty oznacza zapis miejscowy przeglądarki, wartość null oznacza pracę bez zapisu trwałego. */
export function utworzWarstweAod(
  kanal: Kanal,
  opis: OpisKolumnyAod = {},
  magazyn?: MagazynWyciszen | null,
): WarstwaAod {
  const obecnosc = utworzStanObecnosci(magazyn === undefined ? {} : { magazyn });
  const teraz = (): number => Date.now();

  const element = document.createElement('div');
  element.className = 'ao-warstwa';

  // Kontekst wyciszenia jest jedną kopią: menu przy awatarze i w kolumnie czytają ten sam odczyt.
  const kontekstWyciszenia: KontekstWyciszenia = utworzKontekstWyciszenia(kanal);

  const kolumna = utworzKolumneAod(kanal, obecnosc, opis, kontekstWyciszenia);

  const dymek = utworzDymekSugestii({
    stery: kolumna.stery,
    naOdloz: () => {
      // „Odłóż": pozycja wraca do listy oczekujących, status pozostaje `nowa`.
      dymek.schowaj();
      odswiezWyglad();
    },
    naOdrzuc: (wpis) => {
      kolumna.kolejka.odrzuc(wpis.decyzja.klucz);
      kolumna.przerysujDecyzje();
      dymek.schowaj();
      odswiezWyglad();
    },
    naRozwin: () => {
      dymek.schowaj();
      kolumna.ogniskujSugestie();
      odswiezWyglad();
    },
    naZamkniecie: () => {
      dymek.schowaj();
      odswiezWyglad();
    },
  });

  const awatar = utworzAwatarAod({
    naKlikniecie: () => przelaczDymek(),
    naKlikniecePodwojne: () => {
      dymek.schowaj();
      kolumna.przelacz();
      odswiezWyglad();
    },
  });

  // Menu kebab na awatarze to ten sam komponent co w kolumnie — osobny wyzwalacz wyciszenia od ręki.
  const menuAwatara = utworzMenuWyciszenia({
    stan: obecnosc,
    teraz,
    kontekst: kontekstWyciszenia,
    naOtwarcie: () => {
      void kontekstWyciszenia.odswiez().then(() => menuAwatara.odswiez(teraz()));
    },
    etykietaZnaku: 'Menu wyciszania Always On Display — wyciszenie od ręki z pozycji awatara',
  });
  menuAwatara.element.classList.add('ao-wyciszenie--awatar');

  element.append(kolumna.element, dymek.element, awatar.element, menuAwatara.element);

  // Klucz ostatnio ujawnionej decyzji chroni przed powtórnym otwarciem dymka przy tej samej sugestii.
  let ostatnioUjawniona: string | null = null;

  /** Najpilniejszy wpis kolejki albo `null`, gdy nic nie czeka. */
  function najpilniejszy(): WpisKolejkiDecyzji | null {
    const wpisy = kolumna.kolejka.wykaz(teraz());
    if (wpisy.length === 0) return null;
    // Wykaz idzie od najdłużej czekającego; wagę wysoką przepuszcza się przed nią.
    const wysoka = wpisy.find((wpis) => wpis.decyzja.wagaUjawnienia === WagaUjawnienia.Wysoka);
    return wysoka ?? wpisy[0] ?? null;
  }

  function pokazDymek(wpis: WpisKolejkiDecyzji): void {
    dymek.ustawWaski(globalThis.innerWidth - kolumna.szerokosc() < PROG_CZYTELNOSCI);
    dymek.pokaz(wpis, kolumna.dopasujKolejke(wpis.decyzja.idOkna, wpis.decyzja.idSesji));
    odswiezWyglad();
  }

  function przelaczDymek(): void {
    if (dymek.czyWidoczny()) {
      dymek.schowaj();
      odswiezWyglad();
      return;
    }

    const wpis = najpilniejszy();
    if (wpis === null) {
      // Przy pustej kolejce kliknięcie awatara otwiera powierzchnię interakcji, nie pusty dymek.
      kolumna.otworz();
      odswiezWyglad();
      return;
    }

    pokazDymek(wpis);
  }

  // Dymek otwiera się sam tylko przy wadze wysokiej, w progu godzinowym, poza trybem cichym.
  function rozwazUjawnienie(): void {
    const wpis = najpilniejszy();
    if (wpis === null) return;
    if (wpis.decyzja.klucz === ostatnioUjawniona) return;
    if (dymek.czyWidoczny()) return;

    const chwila = teraz();
    if (!obecnosc.czyOtworzycDymek(opisUjawnienia(wpis), chwila)) {
      return;
    }

    ostatnioUjawniona = wpis.decyzja.klucz;
    obecnosc.odnotujDymek(chwila);
    pokazDymek(wpis);
  }

  // Klasa zdarzenia, sesja i moduł opisują decyzję regule wyciszenia; moduł nieznany jest pomijany.
  function opisUjawnienia(wpis: WpisKolejkiDecyzji): OpisUjawnienia {
    const idModulu = kontekstWyciszenia.modulOkna(wpis.decyzja.idOkna);
    return {
      waga: wpis.decyzja.wagaUjawnienia,
      krytyczna: wpis.decyzja.krytyczna,
      klasa: KLASA_POWODU[wpis.decyzja.powod],
      ...(idModulu === undefined ? {} : { idModulu }),
      ...(wpis.decyzja.idSesji === undefined ? {} : { idSesji: wpis.decyzja.idSesji }),
    };
  }

  // Stan awatara wynika ze stanu obecności i z kolejki oczekujących decyzji.
  function stanAwatara(liczba: number, wysoka: boolean, chwila: number): StanAwatara {
    if (obecnosc.tryb() === TrybObecnosci.Ukryty) return StanAwatara.Ukryty;
    // Stan „Wyciszony” obejmuje wyciszenie czasowe, kontekstowe i wyciszenie klasy zdarzeń.
    if (obecnosc.czyJakiekolwiekWyciszenie(chwila)) return StanAwatara.Wyciszony;
    if (liczba === 0) return StanAwatara.Spoczynek;
    return wysoka ? StanAwatara.WagaWysoka : StanAwatara.SugestiaOczekujaca;
  }

  function odswiezWyglad(): void {
    const chwila = teraz();
    const wpisy = kolumna.kolejka.wykaz(chwila);

    // Plakietka liczy sugestie ujawniane, nie wstrzymane wyciszeniem wybiórczym; wyciszenie znaczy awatar.
    const ujawniane = wpisy.filter((wpis) => !obecnosc.czySugestiaWstrzymana(opisUjawnienia(wpis), chwila));
    const wysoka = ujawniane.some((wpis) => wpis.decyzja.wagaUjawnienia === WagaUjawnienia.Wysoka);

    awatar.ustawStan(stanAwatara(ujawniane.length, wysoka, chwila));

    // Wyciszenie czasowe chowa plakietkę, poza sugestią krytyczną, ujawnianą mimo każdego wyciszenia.
    const krytyczne = wpisy.filter((wpis) => wpis.decyzja.krytyczna);
    const plakietkaOgolna = obecnosc.czyPlakietkaWidoczna(chwila);
    awatar.ustawLiczbe(plakietkaOgolna ? ujawniane.length : krytyczne.length);
    awatar.ustawPlakietkeWidoczna(plakietkaOgolna || krytyczne.length > 0);
    menuAwatara.odswiez(chwila);

    const odsuniecie = kolumna.szerokosc();
    awatar.ustawOdsuniecie(odsuniecie);
    // Kebab awatara jedzie z awatarem: stoi pod nim, przy tej samej krawędzi.
    menuAwatara.element.style.setProperty('--ao-odsuniecie', `${Math.max(0, odsuniecie)}px`);
    dymek.ustawOdsuniecie(odsuniecie);
    dymek.ustawWaski(globalThis.innerWidth - odsuniecie < PROG_CZYTELNOSCI);

    element.dataset['kolumna'] = kolumna.czyOtwarta() ? 'otwarta' : 'zamknieta';
    zwezObszarRoboczy(odsuniecie);
  }

  // Zwężenie obszaru roboczego idzie stylem rodzica warstwy, nie zmianą powłoki, więc znika bez śladu.
  function zwezObszarRoboczy(szerokosc: number): void {
    const gospodarz = element.parentElement;
    if (gospodarz === null) return;

    if (szerokosc <= 0) {
      gospodarz.classList.remove('ao-zwezenie');
      gospodarz.style.removeProperty('--ao-zwezenie');
      return;
    }

    gospodarz.classList.add('ao-zwezenie');
    gospodarz.style.setProperty('--ao-zwezenie', `${szerokosc}px`);
  }

  // --- skróty klawiszowe ---------------------------------------------------

  function naKlawisz(zdarzenie: KeyboardEvent): void {
    const modyfikator = (zdarzenie.ctrlKey || zdarzenie.metaKey) && zdarzenie.shiftKey;
    if (!modyfikator) return;

    switch (zdarzenie.key.toLowerCase()) {
      case 'a':
        kolumna.przelacz();
        break;
      case 'v':
        kolumna.ogniskujGlos();
        break;
      case 'q':
        kolumna.ogniskujSugestie();
        break;
      case 'm':
        obecnosc.przelaczWyciszenieKwadransem(teraz());
        break;
      case 'h':
        obecnosc.przelaczUkrycie();
        break;
      default:
        return;
    }

    zdarzenie.preventDefault();
    odswiezWyglad();
  }

  document.addEventListener('keydown', naKlawisz);

  const przestanObserwowacObecnosc = obecnosc.obserwuj(() => {
    if (obecnosc.tryb() === TrybObecnosci.Ukryty) dymek.schowaj();
    odswiezWyglad();
  });

  const przestanObserwowacUklad = kolumna.obserwujUklad(() => {
    odswiezWyglad();
    rozwazUjawnienie();
  });

  odswiezWyglad();

  // Odczyt nadrabiający rusza od razu, bez otwierania kolumny, aby plakietka mówiła prawdę wcześniej.
  void kolumna.odswiezDecyzje().then(() => {
    odswiezWyglad();
    rozwazUjawnienie();
  });

  // Kontekst wyciszenia odświeża się tym torem, by menu nazywało moduł i sesję zanim ktoś je otworzy.
  void kontekstWyciszenia.odswiez().then(() => odswiezWyglad());

  return {
    element,

    otworzPowierzchnie() {
      kolumna.otworz();
      odswiezWyglad();
    },

    obecnosc,
    kolumna,

    rozlacz() {
      document.removeEventListener('keydown', naKlawisz);
      przestanObserwowacObecnosc();
      przestanObserwowacUklad();
      zwezObszarRoboczy(0);
      menuAwatara.rozlacz();
      awatar.rozlacz();
      dymek.rozlacz();
      kolumna.rozlacz();
      element.remove();
    },
  };
}
