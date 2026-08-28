import { Command } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  przelacznikWidoku,
  przestaw,
  przyciskAkcji,
} from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import { oznaczWarstwy, type CzynnoscOkna } from './czynnosci-okna';
import { odtwarzaczSesji, type Odtwarzacz } from './odtwarzanie';
import type { StanTerminala } from './stan-terminala';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import { rysujWiersze, type NastawyWidoku } from './widok-wyjscia';
import type { ZrodloTerminala } from './zrodlo-terminala';

/**
 * Output Console — okno monitorujące modułu Terminal: wynik poleceń ze wszystkich otwartych
 * kart, na żywo, z historią przewijania.
 */
export interface OknoKonsoli {
  element: HTMLElement;
  odswiez(): void;
  /** Czynności okna oddane palecie poleceń i skrótom klawiszowym. */
  czynnosci: readonly CzynnoscOkna[];
}

export function utworzOknoKonsoli(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  pokrycie: PokrycieKomend,
): OknoKonsoli {
  const rama = utworzRameOkna({
    tytul: 'Output Console',
    rola: 'monitor',
    przeznaczenie: 'Zbiorcze wyjście poleceń ze wszystkich kart okna, na żywo, z historią przewijania.',
    kod: 'output-console',
    przedrostek: 'dt',
    ogniskowalne: true,
  });
  const tresc = utworzStanTresci();

  const kontrolki = zlozPowierzchnieKonsoli(rama, tresc.element, () => pokaz(), pokrycie);
  const { wzorzec, zawijanie, znaczniki, przewijanie, regularne, wielkosc, tylkoKarta, podzial } =
    kontrolki;
  const { eksport, czyszczenie, ponownie, doDiagnostyki, odtwarzanie } = kontrolki;

  function pokaz(): void {
    const bufor = stan.bufor();
    const wszystkie = bufor.wiersze();
    if (wszystkie.length === 0) {
      tresc.pusto('Żadne polecenie nie oddało jeszcze wyjścia. Uruchom polecenie w karcie terminala.');
      odtwarzanie.ustawZakres(0);
      return;
    }
    odtwarzanie.ustawZakres(wszystkie.length);

    const idKarty = stan.kartaBiezaca()?.id ?? '';
    const nastawy = nastawyKonsoli(kontrolki, idKarty);
    const wynik = rysujWiersze(wszystkie, nastawy);
    if (wynik.blad !== '') {
      tresc.blad(`Grep nie przyjął wzorca: ${wynik.blad}`);
      return;
    }
    if (wynik.element.childElementCount === 0) {
      // Pustka po zawężeniu ma kilka przyczyn i okno nazywa tę właściwą, nie zawsze wzorzec grepa.
      tresc.pusto(powodPustegoWidoku(kontrolki, stan, wszystkie.length));
      return;
    }
    const miejsce = tresc.tresc();
    const zawezona = nastawy.karta !== undefined;
    if (podzial.dataset['wlaczony'] !== 'true') {
      miejsce.append(wynik.element);
    } else {
      // Podział pokazuje ten sam bufor w dwóch zawężeniach naraz, obok siebie bez przełączania.
      const drugie = rysujWiersze(
        wszystkie,
        zawezona ? { ...nastawy, karta: undefined } : { ...nastawy, karta: idKarty },
      );
      const tafle = document.createElement('div');
      tafle.className = 'dt-tafle';
      tafle.append(
        tafla(zawezona ? 'Karta bieżąca' : 'Wszystkie karty okna', wynik.element),
        tafla(
          zawezona
            ? 'Wszystkie karty okna'
            : idKarty === ''
              ? 'Karta bieżąca — nie ma ani jednej'
              : 'Karta bieżąca',
          drugie.element,
        ),
      );
      miejsce.append(tafle);
    }
    if (bufor.utracone() > 0) miejsce.prepend(uwagaOOgonieHistorii(bufor.utracone()));
    if (przewijanie.dataset['wlaczony'] === 'true') {
      for (const konsola of miejsce.querySelectorAll('.dt-konsola')) {
        konsola.scrollTop = konsola.scrollHeight;
      }
    }
  }

  /** Stawia ognisko w polu grepa; zawężenie do karty bieżącej wchodzi parametrem. */
  function skupGrep(doKartyBiezacej: boolean): void {
    if ((tylkoKarta.dataset['wlaczony'] === 'true') !== doKartyBiezacej) {
      przestaw(tylkoKarta);
      pokaz();
    }
    wzorzec.focus();
    wzorzec.select();
  }

  for (const kontrolkaWidoku of [zawijanie, znaczniki, regularne, wielkosc, tylkoKarta]) {
    kontrolkaWidoku.addEventListener('click', () => {
      przestaw(kontrolkaWidoku);
      rama.element.dataset['zawijanie'] = zawijanie.dataset['wlaczony'] ?? 'true';
      pokaz();
    });
  }
  przewijanie.addEventListener('click', () => przestaw(przewijanie));
  podzial.addEventListener('click', () => {
    rama.element.dataset['podzial'] = String(przestaw(podzial));
  });
  wzorzec.addEventListener('input', pokaz);
  czyszczenie.addEventListener('click', () => {
    const ile = stan.bufor().wiersze().length;
    stan.bufor().wyczysc();
    pokaz();
    tresc.potwierdzenie(
      ile === 0
        ? 'Bufor widoku był już pusty — nie było czego czyścić.'
        : `Bufor widoku wyczyszczony: zniesiono ${ile} wierszy. Rejestr procesów rdzenia zostaje nietknięty.`,
      ile > 0,
    );
  });
  eksport.addEventListener('click', () => {
    const wiersze = stan.bufor().wiersze();
    if (wiersze.length === 0) {
      // Plik o zerowej długości wygląda tak samo jak wyjście utracone.
      tresc.potwierdzenie('Bufor widoku jest pusty — pliku nie zapisano.', false);
      return;
    }
    const nazwa = `wyjscie-${Date.now()}.txt`;
    pobierzPlik(nazwa, wiersze.map((wiersz) => wiersz.tresc).join('\n'), 'text/plain');
    tresc.potwierdzenie(`Zapisano ${nazwa} — ${wiersze.length} wierszy bufora widoku.`, true);
  });
  ponownie.addEventListener('click', () => powtorzOstatniePolecenieKarty(zrodlo, stan, tresc));
  // Powód fragmentu liczy pokrycie z odczytu wykazu komend rdzenia; fragment zostaje w oknie.
  doDiagnostyki.addEventListener('click', () =>
    tresc.potwierdzenie(
      pokrycie.zdanie(Command.ContextTransfer, 'Przekazanie fragmentu do Diagnostics Center'),
      false,
    ),
  );

  stan.naZmiane(pokaz);

  const czynnosci: readonly CzynnoscOkna[] = [
    {
      okno: 'Output Console',
      nazwa: 'Szukaj w karcie bieżącej',
      opis: 'Zawęża strumień do karty bieżącej i stawia ognisko w polu wzorca.',
      warstwa: 'na-zadanie',
      skrot: 'Ctrl/Cmd + F',
      wykonaj: () => skupGrep(true),
    },
    {
      okno: 'Output Console',
      nazwa: 'Szukaj w strumieniu zbiorczym',
      opis: 'Zdejmuje zawężenie do karty i stawia ognisko w polu wzorca.',
      warstwa: 'na-zadanie',
      skrot: 'Ctrl/Cmd + Shift + F',
      wykonaj: () => skupGrep(false),
    },
    {
      okno: 'Output Console',
      nazwa: 'Uruchom ponownie ostatnie polecenie',
      opis: 'Powtarza ostatnie polecenie karty bieżącej w tej samej karcie.',
      warstwa: 'zawsze',
      wykonaj: () => ponownie.click(),
    },
    {
      okno: 'Output Console',
      nazwa: 'Podziel widok strumienia',
      opis: 'Stawia obok siebie strumień zbiorczy i strumień karty bieżącej.',
      warstwa: 'kontekstowa',
      wykonaj: () => podzial.click(),
    },
    {
      okno: 'Output Console',
      nazwa: 'Eksportuj bufor widoku',
      opis: 'Zapisuje wiersze bufora widoku jako plik tekstowy.',
      warstwa: 'kontekstowa',
      wykonaj: () => eksport.click(),
    },
    {
      okno: 'Output Console',
      nazwa: 'Wyczyść bufor widoku',
      opis: 'Znosi wiersze z bufora okna; rejestr procesów rdzenia zostaje nietknięty.',
      warstwa: 'kontekstowa',
      wykonaj: () => czyszczenie.click(),
    },
  ];

  return { element: rama.element, odswiez: pokaz, czynnosci };
}

/** Jedna tafla podzielonego widoku konsoli wraz z podpisem, po którym poznać jej zawężenie karty bieżącej. */
function tafla(podpis: string, konsola: HTMLElement): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-tafla';

  const naglowek = document.createElement('h4');
  naglowek.className = 'dt-tafla__podpis';
  naglowek.textContent = podpis;

  blok.append(naglowek, konsola);
  return blok;
}

/** Nastawy widoku konsoli odczytane wprost z kontrolek; zawężenie karty tylko przy włączonym przełączniku. */
function nastawyKonsoli(kontrolki: PowierzchniaKonsoli, kartaBiezaca: string): NastawyWidoku {
  return {
    wzorzec: kontrolki.wzorzec.value,
    regularne: kontrolki.regularne.dataset['wlaczony'] === 'true',
    wielkoscLiter: kontrolki.wielkosc.dataset['wlaczony'] === 'true',
    znaczniki: kontrolki.znaczniki.dataset['wlaczony'] === 'true',
    ...(kontrolki.tylkoKarta.dataset['wlaczony'] === 'true'
      ? { karta: kartaBiezaca }
      : {}),
    doWiersza: kontrolki.odtwarzanie.polozenie(),
  };
}

/**
 * Powód, dla którego po zawężeniu nie został ani jeden wiersz; wyciąć wszystko potrafią trzy
 * nastawy, a okno wymienia te, które są włączone.
 */
function powodPustegoWidoku(
  kontrolki: PowierzchniaKonsoli,
  stan: StanTerminala,
  wszystkich: number,
): string {
  const powody: string[] = [];
  if (kontrolki.wzorzec.value !== '') {
    powody.push(`wzorzec grepa „${kontrolki.wzorzec.value}"`);
  }
  if (kontrolki.tylkoKarta.dataset['wlaczony'] === 'true') {
    const karta = stan.kartaBiezaca();
    powody.push(
      karta === null
        ? 'zawężenie do karty bieżącej, której nie ma ani jednej'
        : `zawężenie do karty ${karta.title ?? karta.shell}`,
    );
  }
  const polozenie = kontrolki.odtwarzanie.polozenie();
  if (polozenie >= 0) powody.push(`odtwarzanie ustawione na wiersz ${polozenie} z ${wszystkich}`);
  if (powody.length === 0) {
    return `Bufor ma ${wszystkich} wierszy, a mimo to nie został ani jeden — to usterka widoku, nie nastaw.`;
  }
  return `Z ${wszystkich} wierszy bufora nie został ani jeden. Wycina je: ${powody.join(' · ')}.`;
}

/** Uwaga o wierszach zniesionych przez pojemność bufora konsoli — pokazuje się wtedy sam ogon historii. */
function uwagaOOgonieHistorii(utracone: number): HTMLElement {
  const uwaga = document.createElement('p');
  uwaga.className = 'dn-pole-opis';
  uwaga.textContent = `Bufor zniósł ${utracone} najstarszych wierszy — konsola pokazuje ogon historii.`;
  return uwaga;
}

/** Ponowne uruchomienie ostatniego polecenia karty bieżącej; bierze źródło, stan i treść okna parametrem. */
function powtorzOstatniePolecenieKarty(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  tresc: StanTresci,
): void {
  const karta = stan.kartaBiezaca();
  if (karta === null) {
    tresc.blad('Nie ma karty bieżącej — nie ma czego powtórzyć.');
    return;
  }
  const polecenie = stan.ostatniePolecenie(karta.id);
  if (polecenie === '') {
    tresc.blad('Karta bieżąca nie wykonała jeszcze żadnego polecenia.');
    return;
  }
  tresc.ladowanie(`Ponowne uruchomienie „${polecenie}”…`);
  void zrodlo.wykonaj({ sessionId: karta.id, command: polecenie }).then((wynik) => {
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(`Rdzeń nie powtórzył polecenia (${Command.TerminalCommandExec}).`, wynik.blad);
      return;
    }
    stan.zapiszProces(wynik.wynik);
    tresc.potwierdzenie(`Proces ${wynik.wynik.id} uruchomiony ponownie.`, true);
  });
}

/** Kontrolki okna Output Console wraz z odtwarzaczem sesji przewijania historii wierszy jego treści wyjścia. */
interface PowierzchniaKonsoli {
  wzorzec: HTMLInputElement;
  zawijanie: HTMLButtonElement;
  znaczniki: HTMLButtonElement;
  przewijanie: HTMLButtonElement;
  regularne: HTMLButtonElement;
  wielkosc: HTMLButtonElement;
  tylkoKarta: HTMLButtonElement;
  podzial: HTMLButtonElement;
  eksport: HTMLButtonElement;
  czyszczenie: HTMLButtonElement;
  ponownie: HTMLButtonElement;
  doDiagnostyki: HTMLButtonElement;
  odtwarzanie: Odtwarzacz;
}

/**
 * Składa kontrolki, pasek akcji, pasek narzędzi i ciało okna; nie domyka się na stanie okna
 * ani na buforze.
 */
function zlozPowierzchnieKonsoli(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  przyOdtwarzaniu: () => void,
  pokrycie: PokrycieKomend,
): PowierzchniaKonsoli {
  const wzorzec = pole('Grep — wzorzec', 'wzorzec albo wyrażenie regularne');
  const zawijanie = przelacznikWidoku('Zawijaj wiersze', true);
  const znaczniki = przelacznikWidoku('Znaczniki czasu', false);
  const przewijanie = przelacznikWidoku('Auto-przewijanie', true);
  const regularne = przelacznikWidoku('Wyrażenie regularne', false);
  const wielkosc = przelacznikWidoku('Rozróżniaj wielkość liter', false);
  const tylkoKarta = przelacznikWidoku('Tylko karta bieżąca', false);
  const podzial = przelacznikWidoku('Podziel widok', false);

  const eksport = przyciskAkcji('Eksportuj');
  const czyszczenie = przyciskAkcji('Wyczyść');
  const ponownie = przyciskAkcji('Uruchom ponownie ostatnie polecenie', 'dn-btn dn-btn--atrament');
  const doDiagnostyki = przyciskAkcji('Przekaż do Diagnostics');

  const odtwarzanie = odtwarzaczSesji(przyOdtwarzaniu);
  const wyjasnienie = pokrycie.przycisk(
    'Wyjaśnij zaznaczenie modelem',
    Command.MessageSend,
    'Wyjaśnienie zaznaczonego fragmentu; należy ono do okna rozmowy modułu, którego to złożenie nie osadza',
  );

  oznaczWarstwy([
    [wzorzec, 'zawsze'],
    [ponownie, 'zawsze'],
    [tylkoKarta, 'na-zadanie'],
    [regularne, 'na-zadanie'],
    [wielkosc, 'na-zadanie'],
    [zawijanie, 'na-zadanie'],
    [znaczniki, 'kontekstowa'],
    [przewijanie, 'kontekstowa'],
    [podzial, 'kontekstowa'],
    [eksport, 'kontekstowa'],
    [czyszczenie, 'kontekstowa'],
    [doDiagnostyki, 'kontekstowa'],
    [odtwarzanie.element, 'ekspercka'],
    [wyjasnienie, 'ekspercka'],
  ]);

  rama.akcje.append(doDiagnostyki, eksport, czyszczenie, ponownie, podzial);
  rama.narzedzia.append(wzorzec, regularne, wielkosc, tylkoKarta, zawijanie, znaczniki, przewijanie);
  rama.cialo.append(stanTresci);

  rama.akcje.append(odtwarzanie.element, wyjasnienie);

  return {
    wzorzec,
    zawijanie,
    znaczniki,
    przewijanie,
    regularne,
    wielkosc,
    tylkoKarta,
    podzial,
    eksport,
    czyszczenie,
    ponownie,
    doDiagnostyki,
    odtwarzanie,
  };
}
