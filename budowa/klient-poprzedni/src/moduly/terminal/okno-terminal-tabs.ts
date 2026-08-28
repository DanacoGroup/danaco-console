import {
  Command,
  TerminalProcessStatus,
  TerminalSessionStatus,
  TerminalShell,
  type ProcessInitiator,
  type TerminalProcess,
  type TerminalSession,
  type TerminalSessionOpenRequest,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleLiczbowe,
  przelacznikWidoku,
  przestaw,
  przyciskAkcji,
} from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import { oznaczWarstwy, type CzynnoscOkna } from './czynnosci-okna';
import { pasekKart } from './pasek-kart';
import { INICJATORZY_KARTY, opisKarty, POWLOKI, SCHEMATY } from './profil-karty';
import type { StanTerminala } from './stan-terminala';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import { rysujWiersze } from './widok-wyjscia';
import { utworzWyborDrzewem, type WyborDrzewem } from './wybor-drzewem';
import type { ZrodloTerminala } from './zrodlo-terminala';
import { notaZaleznosci, PROGRAMY_POWLOK } from './zaleznosci-zewnetrzne';

/**
 * Terminal Tabs to okno wiodące modułu Terminal: otwiera karty powłok i przełącza między nimi,
 * każda karta jest odrębną sesją własnego rodzaju, katalogu i środowiska.
 */
export interface OknoKart {
  element: HTMLElement;
  odswiez(): void;
  /** Czynności okna oddane palecie poleceń i skrótom klawiszowym. */
  czynnosci: readonly CzynnoscOkna[];
}

export function utworzOknoKart(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  pokrycie: PokrycieKomend,
): OknoKart {
  const rama = utworzRameOkna({
    tytul: 'Terminal Tabs',
    rola: 'wiodące',
    przeznaczenie:
      'Karty powłok okna. Każda karta jest odrębną sesją powłoki: własny rodzaj, katalog i środowisko.',
    kod: 'terminal-tabs',
    przedrostek: 'dt',
    ogniskowalne: true,
  });
  const tresc = utworzStanTresci();
  const karty = pasekKart(stan, {
    wybierz: (idKarty) => stan.ustawKarte(idKarty),
    zamknij: (idKarty) => zamknijKarte(idKarty),
  });

  const kontrolki = zlozPowierzchnieKart(rama, karty.element, tresc.element, pokrycie);

  function pokaz(): void {
    const biezaca = stan.kartaBiezaca();
    karty.odswiez();
    if (biezaca === null) {
      tresc.pusto('Nie ma ani jednej karty. Wybierz powłokę i otwórz nową kartę.');
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.append(opisKarty(biezaca, stan), ogonWyjsciaKarty(biezaca, stan));
  }

  function otworz(zrodlowa: TerminalSession | null): void {
    const okno = stan.okno();
    if (okno === '') {
      tresc.blad('Moduł nie zna okna komunikacji — karty nie ma gdzie otworzyć.');
      return;
    }
    const zadanie = zadanieOtwarciaKarty(okno, zrodlowa, kontrolki);
    otworzKarteRdzenia(zrodlo, stan, tresc, zadanie);
  }

  // Zamknięcie karty najpierw w rdzeniu, potem w widoku; procesy karty zostają biegnące.
  function zamknijKarte(idKarty: string, zProcesami = false): void {
    if (stan.czyPrzypieta(idKarty)) {
      tresc.potwierdzenie('Karta jest przypięta — najpierw ją odepnij. Nic nie zamknięto.', false);
      return;
    }
    tresc.ladowanie(`Zamykanie karty ${idKarty} w rdzeniu…`);
    void zrodlo.zamknijKarte({ sessionId: idKarty, force: zProcesami }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(
          `Rdzeń nie zamknął karty ${idKarty} (${Command.TerminalSessionClose}). ` +
            'Karta zostaje na ekranie, bo powłoka po stronie rdzenia nadal biegnie.',
          wynik.blad,
        );
        return;
      }
      stan.zamknijKarte(idKarty);
      const zatrzymane = wynik.wynik.stoppedProcessIds ?? [];
      tresc.potwierdzenie(
        `Karta ${idKarty} zamknięta w rdzeniu (stan ${wynik.wynik.session.status}). ` +
          (zatrzymane.length === 0
            ? 'Procesy karty biegną dalej — kończy je osobna czynność.'
            : `Zatrzymano wraz z kartą ${zatrzymane.length} procesów.`),
        true,
      );
    });
  }

  /** Czyta z rdzenia karty tego okna — także te sprzed rozłączenia klienta. */
  function odczytajKarty(): void {
    const okno = stan.okno();
    if (okno === '') return;
    void zrodlo.karty({ windowId: okno }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) return;
      for (const karta of wynik.wynik) stan.dodajKarte(karta);
      pokaz();
    });
  }

  /** Czyta plik z katalogu roboczego karty bieżącej i pokazuje jego treść. */
  function odczytajPlikKarty(): void {
    const biezaca = stan.kartaBiezaca();
    if (biezaca === null) {
      tresc.potwierdzenie('Nie ma karty bieżącej — odczyt pliku nie ma czego dotyczyć.', false);
      return;
    }
    const sciezka = kontrolki.plikDoOdczytu.value.trim();
    if (sciezka === '') {
      tresc.potwierdzenie('Odczyt pliku wymaga wskazania ścieżki — wypełnij pole ścieżki.', false);
      return;
    }
    tresc.ladowanie(`Odczyt pliku ${sciezka} na maszynie karty…`);
    void zrodlo.odczytajPlik({ sessionId: biezaca.id, path: sciezka }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie oddał treści pliku ${sciezka} (${Command.TerminalFileRead}).`, wynik.blad);
        return;
      }
      const odpowiedz = wynik.wynik;
      const miejsce = tresc.tresc();
      const podpis = document.createElement('h4');
      podpis.className = 'dt-plik__podpis';
      podpis.textContent = `Treść pliku ${odpowiedz.path}`;
      const tresciwy = document.createElement('pre');
      tresciwy.className = 'dt-plik__tresc';
      tresciwy.textContent = odpowiedz.content;
      miejsce.append(podpis, tresciwy);
      tresc.potwierdzenie(
        odpowiedz.truncated
          ? `Treść przycięta — odcięto ${odpowiedz.truncatedBytes ?? 0} bajtów z początku pliku.`
          : `Plik odczytany w całości${odpowiedz.sizeBytes === undefined ? '' : ` (${odpowiedz.sizeBytes} bajtów)`}.`,
        true,
      );
    });
  }

  kontrolki.nowa.addEventListener('click', () => otworz(null));
  kontrolki.duplikuj.addEventListener('click', () => {
    // Duplikat bez karty źródłowej byłby zwykłą nową kartą, nie duplikatem.
    const biezaca = stan.kartaBiezaca();
    if (biezaca === null) {
      tresc.potwierdzenie('Nie ma karty bieżącej — nie ma czego duplikować.', false);
      return;
    }
    otworz(biezaca);
  });
  kontrolki.uruchom.addEventListener('click', () =>
    wyslijPolecenieKarty(zrodlo, stan, tresc, kontrolki, pokaz),
  );
  // Enter uruchamia polecenie tą samą drogą co przycisk, dla identycznej odpowiedzi.
  kontrolki.polecenie.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter') return;
    if (zdarzenie.shiftKey || zdarzenie.ctrlKey || zdarzenie.altKey || zdarzenie.metaKey) return;
    zdarzenie.preventDefault();
    wyslijPolecenieKarty(zrodlo, stan, tresc, kontrolki, pokaz);
  });
  kontrolki.wstaw.addEventListener('click', () => wstawDoPolaKarty(stan, tresc, kontrolki));
  zwiazCzynnosciKarty(kontrolki, stan, tresc, zamknijKarte, odczytajPlikKarty);
  zwiazCzynnosciWidokuKart(rama.element, kontrolki);

  // Wędrówka po kartach zawija się na końcu wykazu, przy jednej karcie zostaje na miejscu.
  function nastepnaKarta(): void {
    const karty = stan.karty();
    if (karty.length === 0) {
      tresc.potwierdzenie('Nie ma ani jednej karty — nie ma na co przestawić ogniska.', false);
      return;
    }
    if (karty.length === 1) {
      tresc.potwierdzenie('To jedyna karta okna — ognisko zostaje na niej.', false);
      return;
    }
    const biezaca = stan.kartaBiezaca();
    const miejsce = karty.findIndex((karta) => karta.id === biezaca?.id);
    const nastepna = karty[(miejsce + 1) % karty.length];
    if (nastepna === undefined) return;
    stan.ustawKarte(nastepna.id);
    tresc.potwierdzenie(`Ognisko na karcie ${nastepna.title ?? nastepna.shell}.`, true);
  }

  stan.naZmiane(pokaz);

  const czynnosci: readonly CzynnoscOkna[] = [
    {
      okno: 'Terminal Tabs',
      nazwa: 'Otwórz nową kartę powłoki',
      opis: 'Zakłada w rdzeniu kartę powłoki wybranego rodzaju wraz z katalogiem roboczym z formularza.',
      warstwa: 'zawsze',
      skrot: 'Ctrl/Cmd + T',
      wykonaj: () => kontrolki.nowa.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Uruchom polecenie w karcie bieżącej',
      opis: 'Wysyła treść pola polecenia do rdzenia i zapisuje proces w rejestrze.',
      warstwa: 'zawsze',
      wykonaj: () => kontrolki.uruchom.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Zamknij kartę bieżącą',
      opis: 'Zdejmuje kartę z widoku; procesy już uruchomione biegną dalej w rdzeniu.',
      warstwa: 'na-zadanie',
      skrot: 'Ctrl/Cmd + W',
      wykonaj: () => kontrolki.zamknij.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Zamknij kartę wraz z procesami',
      opis: 'Zamyka kartę w rdzeniu i kończy jej procesy sygnałem wymuszonym.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.zamknijZProcesami.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Odczytaj plik karty',
      opis: 'Czyta plik na maszynie karty rdzeniem, bez polecenia powłoki.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.odczytajPlik.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Przestaw ognisko na kolejną kartę',
      opis: 'Przechodzi do następnej karty paska, zawijając się na końcu wykazu.',
      warstwa: 'na-zadanie',
      skrot: 'Ctrl/Cmd + Tab',
      wykonaj: nastepnaKarta,
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Duplikuj kartę bieżącą',
      opis: 'Otwiera kartę o tym samym rodzaju powłoki i katalogu roboczym co bieżąca.',
      warstwa: 'na-zadanie',
      wykonaj: () => kontrolki.duplikuj.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Przypnij albo odepnij kartę',
      opis: 'Karta przypięta stoi przed pozostałymi i nie zamyka się jednym kliknięciem.',
      warstwa: 'na-zadanie',
      wykonaj: () => kontrolki.przypnij.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Zmień nazwę karty',
      opis: 'Nadaje karcie nazwę z pola formularza; nazwa żyje w widoku, rdzeń o niej nie wie.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.przemianuj.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Podziel widok karty na kolumny',
      opis: 'Ustawia opis karty i jej wyjście obok siebie. To układ widoku, nie druga powłoka.',
      warstwa: 'kontekstowa',
      skrot: 'Ctrl/Cmd + Shift + \\',
      wykonaj: () => kontrolki.podzialWidoku.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Eksportuj transkrypt karty',
      opis: 'Zapisuje wiersze bufora widoku należące do karty bieżącej jako plik tekstowy.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.eksport.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Zwiększ pismo konsoli',
      opis: 'Skaluje pismo wszystkich okien modułu w górę.',
      warstwa: 'ekspercka',
      wykonaj: () => kontrolki.wieksze.click(),
    },
    {
      okno: 'Terminal Tabs',
      nazwa: 'Zmniejsz pismo konsoli',
      opis: 'Skaluje pismo wszystkich okien modułu w dół.',
      warstwa: 'ekspercka',
      wykonaj: () => kontrolki.mniejsze.click(),
    },
  ];

  // Odświeżenie czyta karty z rdzenia przed przerysowaniem, ważne po ponownym podłączeniu gniazda.
  function odswiez(): void {
    odczytajKarty();
    pokaz();
  }

  return { element: rama.element, odswiez, czynnosci };
}

/** Ile wierszy ogona bufora wchodzi pod opis karty w widoku — to podgląd wyjścia, nie druga konsola modułu. */
const OGON_KARTY = 120;

/**
 * Ogon wyjścia karty bieżącej pokazuje to, co polecenie wypisało, korzystając z tego samego
 * bufora widoku co konsola, zawężonego do tej karty.
 */
function ogonWyjsciaKarty(karta: TerminalSession, stan: StanTerminala): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-ogon-karty';
  blok.dataset['karta'] = karta.id;

  const podpis = document.createElement('h4');
  podpis.className = 'dt-ogon-karty__podpis';
  podpis.textContent = 'Wyjście tej karty (ogon bufora widoku)';
  blok.append(podpis);

  const wiersze = stan.bufor().wiersze().filter((wiersz) => wiersz.karta === karta.id);
  if (wiersze.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dt-ogon-karty__pusto';
    pusto.textContent =
      'Bufor widoku nie ma ani jednego wiersza tej karty. Albo nic jeszcze nie ' +
      'uruchomiono, albo bufor wyczyszczono, albo wyjście padło przed tym ' +
      'połączeniem — rdzeń może je wtedy nadal mieć, a pokazuje je „Pokaż ' +
      'wyjście" w Process Monitorze.';
    blok.append(pusto);
    return blok;
  }

  // Wzorzec pusty: rysujWiersze nie buduje wyrażenia regularnego, grep należy do konsoli.
  const wynik = rysujWiersze(wiersze, {
    wzorzec: '',
    regularne: false,
    wielkoscLiter: false,
    znaczniki: false,
    karta: karta.id,
    doWiersza: -1,
    oknoRysowania: OGON_KARTY,
  });
  wynik.element.setAttribute('aria-label', `Wyjście karty ${karta.title ?? karta.shell}`);
  blok.append(wynik.element);

  if (wiersze.length > OGON_KARTY) {
    const uwaga = document.createElement('p');
    uwaga.className = 'dn-pole-opis';
    uwaga.textContent =
      `Karta ma ${wiersze.length} wierszy w buforze; tutaj stoi ostatnie ${OGON_KARTY}. ` +
      'Całość jest w Output Console.';
    blok.append(uwaga);
  }
  return blok;
}

/** Ładunek żądania otwarcia karty; puste pole znaczy brak wskazania, więc odpowiedni klucz nie wchodzi wcale do żądania. */
function zadanieOtwarciaKarty(
  okno: string,
  zrodlowa: TerminalSession | null,
  kontrolki: PowierzchniaKart,
): TerminalSessionOpenRequest {
  const wybranaPowloka = (zrodlowa?.shell ?? kontrolki.rodzaj.wartosc()) as TerminalShell;
  const wybranyKatalog = zrodlowa?.workingDir ?? kontrolki.katalog.value.trim();
  const wybranaNazwa =
    zrodlowa === null ? kontrolki.nazwa.value.trim() : `${zrodlowa.title ?? zrodlowa.shell} (kopia)`;
  // Wskazanie celu wchodzi wyłącznie dla powłoki, która go używa, inaczej byłoby bez znaczenia.
  const cel = wskazanieCeluKarty(wybranaPowloka, zrodlowa, kontrolki);
  return {
    windowId: okno,
    shell: wybranaPowloka,
    ...(wybranyKatalog === undefined || wybranyKatalog === '' ? {} : { workingDir: wybranyKatalog }),
    ...(wybranaNazwa === '' ? {} : { title: wybranaNazwa }),
    ...cel,
  };
}

/**
 * Wskazanie celu karty właściwe jej powłoce: kontener, pod, urządzenie albo adres sieciowy,
 * brane z karty źródłowej przy duplikowaniu.
 */
function wskazanieCeluKarty(
  powloka: TerminalShell,
  zrodlowa: TerminalSession | null,
  kontrolki: PowierzchniaKart,
): Partial<TerminalSessionOpenRequest> {
  const cel = zrodlowa?.remoteTarget ?? kontrolki.celKarty.value.trim();
  const port = Number.parseInt(kontrolki.portKarty.value, 10);
  const wskazanie = kontrolki.wskazanieCelu.value.trim();

  switch (powloka) {
    case TerminalShell.Ssh:
    case TerminalShell.Telnet:
      return {
        ...(cel === '' ? {} : { remoteTarget: cel }),
        ...(Number.isFinite(port) && port > 0 ? { remotePort: port } : {}),
      };
    case TerminalShell.Container:
      return wskazanie === '' ? {} : { containerRef: { containerId: wskazanie } };
    case TerminalShell.Pod:
      return wskazanie === '' ? {} : { containerRef: { podName: wskazanie } };
    case TerminalShell.Serial:
      return {
        ...(wskazanie === '' ? {} : { serialDevice: wskazanie }),
        ...(Number.isFinite(port) && port > 0 ? { serialBaudRate: port } : {}),
      };
    default:
      return {};
  }
}

/** Otwarcie karty powłoki w rdzeniu na podstawie złożonego żądania; źródło danych i stan treści wchodzą parametrem. */
function otworzKarteRdzenia(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  tresc: StanTresci,
  zadanie: TerminalSessionOpenRequest,
): void {
  tresc.ladowanie(`Otwieranie karty powłoki ${zadanie.shell}…`);
  void zrodlo.otworzKarte(zadanie).then((wynik) => {
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(`Rdzeń nie otworzył karty powłoki (${Command.TerminalSessionOpen}).`, wynik.blad);
      return;
    }
    const karta = wynik.wynik;
    stan.dodajKarte(karta);
    tresc.potwierdzenie(zdanieOtwarcia(zadanie, karta), karta.status === TerminalSessionStatus.Running);
  });
}

/**
 * Zdanie o otwarciu karty złożone z odpowiedzi rdzenia, nie z żądania, bo to stan karty
 * rozstrzyga o powodzeniu otwarcia.
 */
function zdanieOtwarcia(zadanie: TerminalSessionOpenRequest, karta: TerminalSession): string {
  const oddany = karta.workingDir ?? '';
  const czesci = [
    `Rdzeń otworzył kartę ${karta.shell} (stan karty: ${karta.status}).`,
    `Katalog roboczy wg rdzenia: ${oddany === '' ? 'nie podany' : oddany}.`,
  ];
  const zadany = zadanie.workingDir ?? '';
  if (zadany !== '' && zadany !== oddany) {
    czesci.push(`Żądanie wskazywało ${zadany} — rdzeń wziął inny.`);
  }
  const zadanaNazwa = zadanie.title ?? '';
  const oddanaNazwa = karta.title ?? '';
  if (zadanaNazwa !== '' && zadanaNazwa !== oddanaNazwa) {
    czesci.push(`Nazwa wg rdzenia: ${oddanaNazwa === '' ? 'bez nazwy' : oddanaNazwa} — nie ta, o którą poszło żądanie.`);
  }
  return czesci.join(' ');
}

/** Uruchomienie polecenia w karcie bieżącej okna; odświeżenie widoku po wykonaniu wchodzi wywołaniem zwrotnym. */
function wyslijPolecenieKarty(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  tresc: StanTresci,
  kontrolki: PowierzchniaKart,
  poWykonaniu: () => void,
): void {
  const karta = stan.kartaBiezaca();
  if (karta === null) {
    tresc.blad('Polecenie bez karty nie ma powłoki — otwórz kartę.');
    return;
  }
  const trescPolecenia = kontrolki.polecenie.value.trim();
  if (trescPolecenia === '') {
    tresc.blad('Puste polecenie nie ma czego uruchomić.');
    return;
  }
  tresc.ladowanie(`Uruchamianie „${trescPolecenia}” w karcie ${karta.shell}…`);
  void zrodlo
    .wykonaj({
      sessionId: karta.id,
      command: trescPolecenia,
      initiator: kontrolki.inicjator.wartosc() as ProcessInitiator,
    })
    .then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie uruchomił polecenia (${Command.TerminalCommandExec}).`, wynik.blad);
        return;
      }
      const proces = wynik.wynik;
      stan.zapiszProces(proces);
      // Polecenie do powtórzenia bierze się z odpowiedzi rdzenia, nie z pola formularza.
      stan.zapamietajPolecenie(karta.id, proces.command);
      kontrolki.polecenie.value = '';
      poWykonaniu();
      tresc.potwierdzenie(zdanieUruchomienia(proces), proces.status !== TerminalProcessStatus.Failed);
    });
}

/**
 * Zdanie o uruchomieniu procesu; czasownik bierze się ze stanu z odpowiedzi.
 * Rdzeń potrafi przyjąć wywołanie i w tej samej odpowiedzi oddać proces w stanie
 * `failed`, więc orzekanie uruchomienia z samego faktu wysłania żądania byłoby
 * fałszem.
 */
function zdanieUruchomienia(proces: TerminalProcess): string {
  const ogon = `PID ${proces.pid ?? 'nie podany'}, polecenie „${proces.command}"`;
  if (proces.status === TerminalProcessStatus.Failed) {
    return `Rdzeń przyjął wywołanie, ale proces ${proces.id} jest już w stanie ${proces.status} (${ogon}).`;
  }
  return `Rdzeń uruchomił proces ${proces.id} — stan ${proces.status} (${ogon}).`;
}

/**
 * „Wstaw do terminala" — czynność wyłącznie widokowa, więc mówi tylko o tym, co
 * widać w oknie. Okno niczego tu nie wstawia, co najwyżej stawia ognisko w polu
 * polecenia, dlatego zdanie rozróżnia brak karty od pustego pola.
 */
function wstawDoPolaKarty(
  stan: StanTerminala,
  tresc: StanTresci,
  kontrolki: PowierzchniaKart,
): void {
  const karta = stan.kartaBiezaca();
  if (karta === null) {
    tresc.potwierdzenie('Nie ma karty bieżącej — polecenie nie ma dokąd trafić.', false);
    return;
  }
  const trescPolecenia = kontrolki.polecenie.value.trim();
  kontrolki.polecenie.focus();
  if (trescPolecenia === '') {
    tresc.potwierdzenie(
      `Pole polecenia karty ${karta.shell} jest puste — ognisko ustawione, wpisać trzeba samemu.`,
      false,
    );
    return;
  }
  tresc.potwierdzenie(
    `W polu karty ${karta.shell} stoi „${trescPolecenia}" — do rdzenia nic nie poszło, uruchomienie zostaje decyzją Operatora.`,
    true,
  );
}

/**
 * Czynności karty bieżącej: nazwa, przypięcie, zamknięcie i eksport transkryptu; każde
 * kliknięcie odpowiada, także odmową bez karty bieżącej.
 */
function zwiazCzynnosciKarty(
  kontrolki: PowierzchniaKart,
  stan: StanTerminala,
  tresc: StanTresci,
  zamknijKarte: (idKarty: string, zProcesami?: boolean) => void,
  odczytajPlikKarty: () => void,
): void {
  /** Karta bieżąca albo odmowa z powodem — jedno zdanie dla czterech czynności. */
  function karta(czynnosc: string): TerminalSession | null {
    const biezaca = stan.kartaBiezaca();
    if (biezaca === null) {
      tresc.potwierdzenie(`Nie ma karty bieżącej — ${czynnosc} nie ma czego dotyczyć.`, false);
    }
    return biezaca;
  }

  kontrolki.przemianuj.addEventListener('click', () => {
    const biezaca = karta('zmiana nazwy');
    if (biezaca === null) return;
    const nazwa = kontrolki.nazwa.value.trim();
    if (nazwa === '') {
      tresc.potwierdzenie('Pole nazwy jest puste — nazwy karty nie zmieniono.', false);
      return;
    }
    stan.przemianujKarte(biezaca.id, nazwa);
    // Nazwa karty żyje wyłącznie w widoku: kontrakt nie ma komendy zmiany nazwy karty powłoki.
    tresc.potwierdzenie(`Nazwa karty w widoku to teraz „${nazwa}". Rdzeń o zmianie nie wie.`, true);
  });

  kontrolki.przypnij.addEventListener('click', () => {
    const biezaca = karta('przypięcie');
    if (biezaca === null) return;
    const przypieta = przestaw(kontrolki.przypnij);
    stan.przypnijKarte(biezaca.id, przypieta);
    tresc.potwierdzenie(
      przypieta
        ? 'Karta przypięta — stoi przed pozostałymi i nie zamknie się jednym kliknięciem.'
        : 'Karta odpięta — zamknięcie znów działa jednym kliknięciem.',
      true,
    );
  });

  kontrolki.zamknij.addEventListener('click', () => {
    const biezaca = karta('zamknięcie');
    if (biezaca !== null) zamknijKarte(biezaca.id);
  });

  kontrolki.zamknijZProcesami.addEventListener('click', () => {
    const biezaca = karta('zamknięcie wraz z procesami');
    if (biezaca !== null) zamknijKarte(biezaca.id, true);
  });

  kontrolki.odczytajPlik.addEventListener('click', odczytajPlikKarty);

  kontrolki.eksport.addEventListener('click', () => {
    const biezaca = karta('eksport transkryptu');
    if (biezaca === null) return;
    const wiersze = stan.bufor().wiersze().filter((wiersz) => wiersz.karta === biezaca.id);
    if (wiersze.length === 0) {
      // Plik o zerowej długości wygląda tak samo jak transkrypt utracony.
      tresc.potwierdzenie('Ta karta nie ma ani jednego wiersza wyjścia — pliku nie zapisano.', false);
      return;
    }
    const nazwaPliku = `transkrypt-${biezaca.id}.txt`;
    pobierzPlik(nazwaPliku, wiersze.map((wiersz) => wiersz.tresc).join('\n'), 'text/plain');
    tresc.potwierdzenie(`Zapisano ${nazwaPliku} — ${wiersze.length} wierszy z bufora widoku.`, true);
  });
}

/**
 * Czynności wyłącznie widokowe: schemat barw, podział widoku i rozmiar pisma; kontrakt nie ma
 * dla nich komend, więc nie sięgają po rdzeń ani stan modułu.
 */
function zwiazCzynnosciWidokuKart(element: HTMLElement, kontrolki: PowierzchniaKart): void {
  let rozmiarPisma = 100;

  /** Węzeł modułu albo samo okno, gdy złożenie stoi poza modułem. */
  function nosnikNastaw(): HTMLElement {
    return element.closest<HTMLElement>('.dt-modul') ?? element;
  }

  function ustawPismo(zmiana: number): void {
    rozmiarPisma = Math.min(200, Math.max(60, rozmiarPisma + zmiana));
    nosnikNastaw().style.setProperty('--dt-pismo', `${rozmiarPisma}%`);
  }

  kontrolki.schemat.naZmiane((wartosc) => {
    nosnikNastaw().dataset['schemat'] = wartosc;
  });
  kontrolki.podzialWidoku.addEventListener('click', () => {
    element.dataset['podzialWidoku'] = String(przestaw(kontrolki.podzialWidoku));
  });
  kontrolki.mniejsze.addEventListener('click', () => ustawPismo(-10));
  kontrolki.wieksze.addEventListener('click', () => ustawPismo(10));
}

/** Kontrolki formularza i przycisków okna Terminal Tabs: wybór powłoki, katalog, polecenie, karty i widok. */
interface PowierzchniaKart {
  rodzaj: WyborDrzewem;
  katalog: HTMLInputElement;
  nazwa: HTMLInputElement;
  polecenie: HTMLInputElement;
  inicjator: WyborDrzewem;
  schemat: WyborDrzewem;
  nowa: HTMLButtonElement;
  duplikuj: HTMLButtonElement;
  przemianuj: HTMLButtonElement;
  przypnij: HTMLButtonElement;
  zamknij: HTMLButtonElement;
  celKarty: HTMLInputElement;
  wskazanieCelu: HTMLInputElement;
  portKarty: HTMLInputElement;
  zamknijZProcesami: HTMLButtonElement;
  plikDoOdczytu: HTMLInputElement;
  odczytajPlik: HTMLButtonElement;
  uruchom: HTMLButtonElement;
  wstaw: HTMLButtonElement;
  eksport: HTMLButtonElement;
  podzialWidoku: HTMLButtonElement;
  mniejsze: HTMLButtonElement;
  wieksze: HTMLButtonElement;
}

/**
 * Składa kontrolki, pasek akcji, pasek narzędzi i ciało okna; pozycje niewykonywane są
 * nieczynne z powodem liczonym z wykazu komend rdzenia.
 */
function zlozPowierzchnieKart(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  pasek: HTMLElement,
  stanTresci: HTMLElement,
  pokrycie: PokrycieKomend,
): PowierzchniaKart {
  const rodzaj = utworzWyborDrzewem({ nastawa: 'Powłoka nowej karty', pozycje: POWLOKI });
  const katalog = pole('Katalog roboczy karty', 'puste = własny katalog okna');
  const nazwa = pole('Nazwa karty', 'np. budowanie');
  const polecenie = pole('Polecenie powłoki', 'np. npm run build — Enter uruchamia');
  const inicjator = utworzWyborDrzewem({
    nastawa: 'Inicjator uruchomienia',
    pozycje: INICJATORZY_KARTY,
  });
  const schemat = utworzWyborDrzewem({ nastawa: 'Schemat barw kart', pozycje: SCHEMATY });

  const nowa = przyciskAkcji('+ Nowa karta', 'dn-btn dn-btn--atrament');
  const duplikuj = przyciskAkcji('Duplikuj kartę');
  const przemianuj = przyciskAkcji('Zmień nazwę');
  const przypnij = przelacznikWidoku('Przypnij kartę', false);
  const zamknij = przyciskAkcji('Zamknij kartę');
  const uruchom = przyciskAkcji('Uruchom teraz', 'dn-btn dn-btn--atrament');
  const wstaw = przyciskAkcji('Wstaw do terminala');
  const eksport = przyciskAkcji('Eksportuj transkrypt');
  const podzialWidoku = przelacznikWidoku('Podziel widok karty na kolumny', false);
  podzialWidoku.title =
    'Ustawia opis karty i jej wyjście obok siebie. To układ widoku, nie druga powłoka — kontrakt nie zna panelu wewnątrz karty.';
  const mniejsze = przyciskAkcji('A–');
  const wieksze = przyciskAkcji('A+');

  const zamknijZProcesami = przyciskAkcji('Zamknij kartę wraz z procesami');
  zamknijZProcesami.title =
    'Zamyka kartę w rdzeniu i kończy jej procesy sygnałem wymuszonym. Czynność jest osobna, bo ' +
    'zamknięcie zakładki nie ma prawa przerwać budowania, które trwa trzecią minutę.';
  const celKarty = pole(
    'Adres celu karty zdalnej albo Telnet',
    'użytkownik@host albo alias konfiguracji serwera',
  );
  const wskazanieCelu = pole(
    'Kontener, pod albo urządzenie karty',
    'np. identyfikator kontenera, nazwa poda, /dev/ttyUSB0',
  );
  const portKarty = poleLiczbowe(
    'Port celu albo prędkość portu szeregowego',
    'puste = wartość domyślna',
  );
  const plikDoOdczytu = pole(
    'Ścieżka pliku do odczytu',
    'np. package.json — względem katalogu roboczego karty',
  );
  const odczytajPlik = przyciskAkcji('Odczytaj plik karty');
  odczytajPlik.title =
    'Czyta plik na maszynie karty i pokazuje jego treść. Odczyt idzie rdzeniem, nie poleceniem ' +
    'powłoki — nie zależy więc ani od programu wypisującego plik, ani od składni powłoki karty.';
  // Nazwa spoza kontraktu jest wskazaniem, nie zapisem stanu; karta to najmniejsza jednostka rdzenia.
  const podzialPaneli = pokrycie.przycisk(
    'Podziel kartę na panele powłoki',
    'terminal.pane.split',
    'Podział karty na sąsiadujące panele powłoki',
  );

  oznaczWarstwy([
    [nowa, 'zawsze'],
    [uruchom, 'zawsze'],
    [polecenie, 'zawsze'],
    [rodzaj.element, 'na-zadanie'],
    [katalog, 'na-zadanie'],
    [nazwa, 'na-zadanie'],
    [inicjator.element, 'na-zadanie'],
    [wstaw, 'na-zadanie'],
    [duplikuj, 'na-zadanie'],
    [zamknij, 'na-zadanie'],
    [zamknijZProcesami, 'kontekstowa'],
    [celKarty, 'na-zadanie'],
    [wskazanieCelu, 'kontekstowa'],
    [portKarty, 'kontekstowa'],
    [plikDoOdczytu, 'kontekstowa'],
    [odczytajPlik, 'kontekstowa'],
    [przypnij, 'na-zadanie'],
    [przemianuj, 'kontekstowa'],
    [eksport, 'kontekstowa'],
    [podzialWidoku, 'kontekstowa'],
    [schemat.element, 'kontekstowa'],
    [podzialPaneli, 'ekspercka'],
    [mniejsze, 'ekspercka'],
    [wieksze, 'ekspercka'],
  ]);

  rama.akcje.append(
    nowa,
    duplikuj,
    przemianuj,
    przypnij,
    zamknij,
    zamknijZProcesami,
    odczytajPlik,
    eksport,
    podzialWidoku,
    podzialPaneli,
    mniejsze,
    wieksze,
  );
  rama.narzedzia.append(
    rodzaj.element,
    katalog,
    nazwa,
    celKarty,
    wskazanieCelu,
    portKarty,
    schemat.element,
    polecenie,
    inicjator.element,
    uruchom,
    wstaw,
    plikDoOdczytu,
  );
  rama.cialo.append(
    pasek,
    notaZaleznosci(
      'Karta powłoki nie jest emulatorem wbudowanym w platformę: rdzeń uruchamia program powłoki leżący ' +
        'na maszynie serwera. Instalka Danaco Console nie niesie ani jednego z tych programów, a każde ' +
        'polecenie karty jest odrębnym procesem — stan powłoki, w tym katalog zmieniony poleceniem, ' +
        'nie przechodzi do polecenia następnego.',
      PROGRAMY_POWLOK,
    ),
    stanTresci,
  );

  return {
    celKarty,
    wskazanieCelu,
    portKarty,
    zamknijZProcesami,
    plikDoOdczytu,
    odczytajPlik,
    rodzaj,
    katalog,
    nazwa,
    polecenie,
    inicjator,
    schemat,
    nowa,
    duplikuj,
    przemianuj,
    przypnij,
    zamknij,
    uruchom,
    wstaw,
    eksport,
    podzialWidoku,
    mniejsze,
    wieksze,
  };
}
