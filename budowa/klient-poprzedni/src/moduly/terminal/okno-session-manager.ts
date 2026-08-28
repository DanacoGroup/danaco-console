import {
  Command,
  TerminalKeyType,
  TerminalSessionStatus,
  TerminalShell,
  TerminalTunnelKind,
  TerminalTunnelStatus,
  type TerminalHost,
  type TerminalSession,
  type TerminalSessionOpenRequest,
  type TerminalSshKey,
  type TerminalTunnel,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleLiczbowe,
  pozycjaWykazu,
  przyciskAkcji,
  wybor,
  wykaz,
} from '../../modele/kontrolki-formularza';
import { oznaczWarstwy, type CzynnoscOkna } from './czynnosci-okna';
import {
  czytajKonfiguracjeSsh,
  utworzKsiazkeHostow,
  zapiszKonfiguracjeSsh,
  type KsiazkaHostow,
  type WpisHosta,
} from './ksiazka-hostow';
import type { PokrycieKomend } from '../pokrycie-komend';
import type { StanTerminala } from './stan-terminala';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import type { ZrodloTerminala } from './zrodlo-terminala';
import {
  notaZaleznosci,
  programPowloki,
  zdanieProgramuPowloki,
  type ProgramZewnetrzny,
} from './zaleznosci-zewnetrzne';

/**
 * Session Manager to okno modułu Terminal spinające książkę hostów, karty powłok oraz wymianę
 * z plikiem konfiguracyjnym OpenSSH z bytami rdzenia aplikacji.
 */
export interface OknoZarzadcySesji {
  element: HTMLElement;
  odswiez(): void;
  czynnosci: readonly CzynnoscOkna[];
}

/**
 * Nazwa zmiennej środowiska, którą rdzeń czytał jako adres powłoki zdalnej przed wprowadzeniem
 * pola remoteTarget kontraktu; okno wysyła dziś oba wskazania dla zgodności ze starszym rdzeniem.
 */
const ZMIENNA_CELU = 'SSH_TARGET';

export function utworzOknoZarzadcySesji(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  pokrycie: PokrycieKomend,
): OknoZarzadcySesji {
  const rama = utworzRameOkna({
    tytul: 'Session Manager',
    rola: 'zarządca',
    przeznaczenie:
      'Książka hostów, połączenia zdalne i karty powłok tego okna. Adres celu jedzie do rdzenia zmienną środowiska karty.',
    kod: 'session-manager',
    przedrostek: 'dt',
    ogniskowalne: true,
  });
  const tresc = utworzStanTresci();
  const ksiazka = utworzKsiazkeHostow();
  // Klucze, tunele i karty rdzenia trzymane są w polach okna, a nie w książce, bo to trzy różne byty.
  let klucze: readonly TerminalSshKey[] = [];
  let tunele: readonly TerminalTunnel[] = [];
  let kartyRdzeniaWykaz: readonly TerminalSession[] = [];
  const kontrolki = zlozPowierzchnieHostow(rama, tresc.element, pokrycie);

  function pokaz(): void {
    const miejsce = tresc.tresc();
    miejsce.append(
      wykazHostow(ksiazka, {
        polacz: (wpis) => polaczZHostem(zrodlo, stan, tresc, wpis),
        wczytaj: (wpis) => wpiszDoFormularza(kontrolki, wpis),
        usun: (wpis) => usunHosta(wpis),
      }),
      wykazKluczy(klucze, (klucz) => zdejmijKlucz(klucz)),
      wykazTuneli(tunele, (tunel) => zamknijTunel(tunel)),
      wykazKart(kartyRdzeniaWykaz, stan),
    );
  }

  /** Czyta z rdzenia wszystkie cztery wykazy okna. */
  function odczytaj(): void {
    const okno = stan.okno();
    void zrodlo.hosty({}).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie oddał książki hostów (${Command.TerminalHostList}).`, wynik.blad);
        return;
      }
      ksiazka.zastap(wynik.wynik.map(wpisZKontraktu));
      pokaz();
    });
    void zrodlo.klucze().then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) return;
      klucze = wynik.wynik;
      // Wykaz kluczy zasila zarazem listę wyboru przy wpisie hosta: klucz da się wskazać nazwą.
      const wybrany = kontrolki.klucz.value;
      kontrolki.klucz.replaceChildren();
      for (const [wartosc, opis] of [
        ['', 'klucz domyślny konfiguracji maszyny rdzenia'] as const,
        ...klucze.map((klucz) => [klucz.id, `${klucz.name} (${klucz.keyType})`] as const),
      ]) {
        const pozycja = document.createElement('option');
        pozycja.value = wartosc;
        pozycja.textContent = opis;
        kontrolki.klucz.append(pozycja);
      }
      kontrolki.klucz.value = klucze.some((klucz) => klucz.id === wybrany) ? wybrany : '';
      pokaz();
    });
    void zrodlo.tunele(okno === '' ? {} : { windowId: okno }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) return;
      tunele = wynik.wynik;
      pokaz();
    });
    void zrodlo.karty(okno === '' ? {} : { windowId: okno }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) return;
      kartyRdzeniaWykaz = wynik.wynik;
      pokaz();
    });
  }

  /** Zdejmuje wpis z książki RDZENIA, a nie tylko z wykazu na ekranie. */
  function usunHosta(wpis: WpisHosta): void {
    if (wpis.id === undefined || wpis.id === '') {
      // Wpis bez identyfikatora nigdy nie doszedł do rdzenia; znika z samego widoku, bo nie był zapisany.
      ksiazka.usun(wpis.nazwa);
      pokaz();
      tresc.potwierdzenie(
        `Wpis ${wpis.nazwa} zniesiony z wykazu. W rdzeniu go nie było — nie został zapisany.`,
        true,
      );
      return;
    }
    void zrodlo.usunHosta({ hostId: wpis.id }).then((wynik) => {
      if (!wynik.udany) {
        tresc.blad(
          `Rdzeń nie zdjął wpisu ${wpis.nazwa} (${Command.TerminalHostRemove}).`,
          wynik.blad,
        );
        return;
      }
      odczytaj();
      tresc.potwierdzenie(
        wynik.wynik === true
          ? `Wpis ${wpis.nazwa} zdjęty z książki hostów. Karty już otwarte do tego hosta biegną dalej.`
          : `Wpisu ${wpis.nazwa} w książce rdzenia nie było — wykaz odświeżony.`,
        true,
      );
    });
  }

  /** Zdejmuje klucz z wykazu rdzenia; plików na dysku nie tyka. */
  function zdejmijKlucz(klucz: TerminalSshKey): void {
    void zrodlo.usunKlucz({ keyId: klucz.id }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie zdjął klucza ${klucz.name} (${Command.TerminalKeyRemove}).`, wynik.blad);
        return;
      }
      const odpiete = wynik.wynik.detachedHostIds ?? [];
      odczytaj();
      tresc.potwierdzenie(
        `Klucz ${klucz.name} zdjęty z wykazu; pliki na dysku maszyny rdzenia zostały nietknięte. ` +
          (odpiete.length === 0
            ? 'Żaden wpis książki go nie wskazywał.'
            : `Wpisy, które go wskazywały (${odpiete.length}), wróciły do klucza domyślnego konfiguracji maszyny.`),
        wynik.wynik.removed,
      );
    });
  }

  /** Zamyka przekierowanie portu w rdzeniu. */
  function zamknijTunel(tunel: TerminalTunnel): void {
    void zrodlo.zamknijTunel({ tunnelId: tunel.id }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie zamknął tunelu ${tunel.id} (${Command.TerminalTunnelClose}).`, wynik.blad);
        return;
      }
      odczytaj();
      tresc.potwierdzenie(`Tunel ${tunel.id} zamknięty — stan: ${wynik.wynik.status}.`, true);
    });
  }

  kontrolki.dodaj.addEventListener('click', () => {
    const wpis = wpisZFormularza(kontrolki);
    if (wpis === null) {
      tresc.potwierdzenie('Wpis bez nazwy i bez adresu celu nie ma czego otwierać — nie dodano.', false);
      return;
    }
    zapiszHosta(wpis);
  });

  kontrolki.polacz.addEventListener('click', () => {
    const wpis = wpisZFormularza(kontrolki);
    if (wpis === null) {
      tresc.blad('Połączenie bez adresu celu nie ma dokąd pójść — wypełnij pole adresu.');
      return;
    }
    polaczZHostem(zrodlo, stan, tresc, wpis);
  });

  /** Zapisuje wpis w KSIĄŻCE RDZENIA i odświeża wykaz jego odpowiedzią. */
  function zapiszHosta(wpis: WpisHosta): void {
    const host: TerminalHost = {
      id: wpis.id ?? '',
      name: wpis.nazwa,
      target: wpis.cel,
      createdAt: 0,
      updatedAt: 0,
      ...(wpis.port === undefined ? {} : { port: wpis.port }),
      ...(wpis.grupa === '' ? {} : { group: wpis.grupa }),
      ...(wpis.katalog === '' ? {} : { workingDir: wpis.katalog }),
      ...(wpis.notatka === '' ? {} : { note: wpis.notatka }),
      ...(kontrolki.klucz.value === '' ? {} : { keyId: kontrolki.klucz.value }),
    };
    void zrodlo.zapiszHosta({ host }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie zapisał wpisu ${wpis.nazwa} (${Command.TerminalHostSave}).`, wynik.blad);
        return;
      }
      odczytaj();
      tresc.potwierdzenie(
        wynik.wynik.created
          ? `Wpis ${wynik.wynik.host.name} zapisany w książce rdzenia — przeżyje odświeżenie strony i restart rdzenia.`
          : `Wpis ${wynik.wynik.host.name} podmieniony w książce rdzenia.`,
        true,
      );
    });
  }

  kontrolki.zapiszWszystkie.addEventListener('click', () => {
    const wpisy = ksiazka.wpisy().filter((wpis) => wpis.id === undefined || wpis.id === '');
    if (wpisy.length === 0) {
      tresc.potwierdzenie('Każdy wpis wykazu ma już swój wiersz w rdzeniu — nie ma czego zapisywać.', false);
      return;
    }
    for (const wpis of wpisy) zapiszHosta(wpis);
  });

  kontrolki.wytworzKlucz.addEventListener('click', () => {
    const nazwa = kontrolki.nazwaKlucza.value.trim();
    if (nazwa === '') {
      tresc.blad('Klucz bez nazwy nie da się później wskazać — wypełnij pole nazwy klucza.');
      return;
    }
    tresc.ladowanie(`Wytwarzanie klucza ${nazwa} na maszynie rdzenia…`);
    void zrodlo
      .wytworzKlucz({ name: nazwa, keyType: kontrolki.rodzajKlucza.value as TerminalKeyType })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad(`Rdzeń nie wytworzył klucza ${nazwa} (${Command.TerminalKeyGenerate}).`, wynik.blad);
          return;
        }
        odczytaj();
        tresc.potwierdzenie(
          `Klucz ${wynik.wynik.name} wytworzony na maszynie rdzenia. Odcisk: ${wynik.wynik.fingerprint}. ` +
            'Część tajna nie opuściła maszyny rdzenia; na host docelowy przenosi się klucz publiczny.',
          true,
        );
      });
  });

  kontrolki.wciagnijKlucz.addEventListener('click', () => {
    const nazwa = kontrolki.nazwaKlucza.value.trim();
    const sciezka = kontrolki.sciezkaKlucza.value.trim();
    if (nazwa === '' || sciezka === '') {
      tresc.blad('Wciągnięcie klucza wymaga jego nazwy i ścieżki na maszynie rdzenia.');
      return;
    }
    void zrodlo.wciagnijKlucz({ name: nazwa, path: sciezka }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(`Rdzeń nie wciągnął klucza ze ścieżki ${sciezka} (${Command.TerminalKeyImport}).`, wynik.blad);
        return;
      }
      odczytaj();
      tresc.potwierdzenie(
        `Klucz ${wynik.wynik.name} wciągnięty do wykazu (${wynik.wynik.keyType}, odcisk ${wynik.wynik.fingerprint}). ` +
          (wynik.wynik.hasPassphrase ? 'Klucz jest chroniony hasłem — poda je program ssh.' : ''),
        true,
      );
    });
  });

  kontrolki.otworzTunel.addEventListener('click', () => {
    const okno = stan.okno();
    if (okno === '') {
      tresc.blad('Moduł nie zna okna komunikacji — tunelu nie ma gdzie założyć.');
      return;
    }
    const cel = kontrolki.cel.value.trim();
    const portZdalny = Number.parseInt(kontrolki.portZdalny.value, 10);
    const portLokalny = Number.parseInt(kontrolki.portLokalny.value, 10);
    const rodzaj = kontrolki.rodzajTunelu.value as TerminalTunnelKind;
    if (cel === '') {
      tresc.blad('Tunel bez adresu celu nie ma dokąd pójść — wypełnij pole adresu celu.');
      return;
    }
    if (rodzaj !== TerminalTunnelKind.Dynamic && !Number.isFinite(portZdalny)) {
      tresc.blad('Przekierowanie miejscowe i zwrotne wymaga portu docelowego po drugiej stronie tunelu.');
      return;
    }
    tresc.ladowanie(`Zakładanie tunelu ${rodzaj} przez ${cel}…`);
    void zrodlo
      .otworzTunel({
        windowId: okno,
        kind: rodzaj,
        remoteTarget: cel,
        ...(Number.isFinite(portLokalny) ? { localPort: portLokalny } : {}),
        ...(Number.isFinite(portZdalny) ? { remotePort: portZdalny } : {}),
        ...(kontrolki.hostDocelowy.value.trim() === ''
          ? {}
          : { remoteHost: kontrolki.hostDocelowy.value.trim() }),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad(`Rdzeń nie założył tunelu (${Command.TerminalTunnelOpen}).`, wynik.blad);
          return;
        }
        odczytaj();
        const tunel = wynik.wynik;
        tresc.potwierdzenie(
          tunel.status === TerminalTunnelStatus.Active
            ? `Tunel ${tunel.id} stoi — port ${tunel.localPort ?? '?'} maszyny rdzenia prowadzi przez ${cel}.`
            : `Tunel ${tunel.id} nie stanął (${tunel.status}): ${tunel.errorMessage ?? 'rdzeń nie podał powodu'}.`,
          tunel.status === TerminalTunnelStatus.Active,
        );
      });
  });

  kontrolki.plik.addEventListener('change', () => {
    const wybrany = kontrolki.plik.files?.[0];
    if (wybrany === undefined) return;
    void wybrany.text().then((zawartosc) => {
      const wczytane = czytajKonfiguracjeSsh(zawartosc);
      // Nazwa pliku wraca do pustej, żeby dało się wczytać ten sam plik drugi raz z rzędu.
      kontrolki.plik.value = '';
      if (wczytane.length === 0) {
        tresc.potwierdzenie(
          `Plik ${wybrany.name} nie ma ani jednego wpisu Host z adresem — książki nie zmieniono.`,
          false,
        );
        return;
      }
      const bilans = ksiazka.scal(wczytane);
      pokaz();
      tresc.potwierdzenie(
        `Z pliku ${wybrany.name} weszło ${bilans.nowe} wpisów nowych i ${bilans.podmienione} podmienionych. ` +
          'Wpisy stoją na razie w samym wykazie — do książki rdzenia wnosi je pozycja „Zapisz wykaz w rdzeniu”. ' +
          'To plik z maszyny Operatora; adresy rozwiązuje program ssh uruchomiony na serwerze rdzenia.',
        true,
      );
    });
  });

  kontrolki.eksport.addEventListener('click', () => {
    const wpisy = ksiazka.wpisy();
    if (wpisy.length === 0) {
      // Plik z samym nagłówkiem wygląda tak samo jak wywóz, który się nie udał.
      tresc.potwierdzenie('Książka hostów jest pusta — pliku nie zapisano.', false);
      return;
    }
    const nazwa = `ksiazka-hostow-${Date.now()}.ssh-config`;
    pobierzPlik(nazwa, zapiszKonfiguracjeSsh(wpisy), 'text/plain');
    tresc.potwierdzenie(`Zapisano ${nazwa} — ${wpisy.length} wpisów książki widoku.`, true);
  });

  // Okno budzi się wyłącznie na zmianę wykazu kart, nie na każdy fragment wyjścia poleceń.
  let podpisKart = '';
  stan.naZmiane(() => {
    const podpis = stan
      .karty()
      .map((karta) => `${karta.id}:${karta.status}`)
      .join('|');
    if (podpis === podpisKart) return;
    podpisKart = podpis;
    pokaz();
  });

  const czynnosci: readonly CzynnoscOkna[] = [
    {
      okno: 'Session Manager',
      nazwa: 'Połącz z hostem z formularza',
      opis: 'Otwiera kartę powłoki zdalnej dla adresu wpisanego w polach okna.',
      warstwa: 'zawsze',
      wykonaj: () => kontrolki.polacz.click(),
    },
    {
      okno: 'Session Manager',
      nazwa: 'Dodaj host do książki',
      opis: 'Zapisuje adres z formularza jako wpis książki widoku.',
      warstwa: 'na-zadanie',
      wykonaj: () => kontrolki.dodaj.click(),
    },
    {
      okno: 'Session Manager',
      nazwa: 'Wczytaj konfigurację OpenSSH',
      opis: 'Czyta plik ~/.ssh/config z maszyny Operatora i składa z niego wpisy książki.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.plik.click(),
    },
    {
      okno: 'Session Manager',
      nazwa: 'Zapisz wykaz w książce rdzenia',
      opis: 'Wnosi do dziennika rdzenia wpisy, które stoją jeszcze tylko w wykazie okna.',
      warstwa: 'na-zadanie',
      wykonaj: () => kontrolki.zapiszWszystkie.click(),
    },
    {
      okno: 'Session Manager',
      nazwa: 'Wytwórz klucz SSH',
      opis: 'Wytwarza parę kluczy na maszynie rdzenia; część tajna jej nie opuszcza.',
      warstwa: 'ekspercka',
      wykonaj: () => kontrolki.wytworzKlucz.click(),
    },
    {
      okno: 'Session Manager',
      nazwa: 'Wciągnij klucz do wykazu',
      opis: 'Wciąga do wykazu klucz leżący już na maszynie rdzenia — wskazany ścieżką, nie treścią.',
      warstwa: 'ekspercka',
      wykonaj: () => kontrolki.wciagnijKlucz.click(),
    },
    {
      okno: 'Session Manager',
      nazwa: 'Załóż tunel portowy',
      opis: 'Zakłada przekierowanie portu przez adres celu z formularza.',
      warstwa: 'ekspercka',
      wykonaj: () => kontrolki.otworzTunel.click(),
    },
    {
      okno: 'Session Manager',
      nazwa: 'Odczytaj wykazy z rdzenia',
      opis: 'Czyta z rdzenia książkę hostów, wykaz kluczy, tunele i karty powłok.',
      warstwa: 'kontekstowa',
      wykonaj: () => odczytaj(),
    },
    {
      okno: 'Session Manager',
      nazwa: 'Eksportuj książkę hostów',
      opis: 'Zapisuje książkę widoku jako plik w składni konfiguracji OpenSSH.',
      warstwa: 'kontekstowa',
      wykonaj: () => kontrolki.eksport.click(),
    },
  ];

  return { element: rama.element, odswiez: odczytaj, czynnosci };
}

/**
 * Czynności wiersza hosta: połączenie, wczytanie do formularza i usunięcie; okno oddaje je
 * wykazowi, bo wykaz nie zna ani rdzenia, ani stanu okna.
 */
interface CzynnosciHosta {
  polacz(wpis: WpisHosta): void;
  wczytaj(wpis: WpisHosta): void;
  usun(wpis: WpisHosta): void;
}

/**
 * Wykaz hostów pogrupowany po folderze książki; pustka wykazu ma własne zdanie ekranu, bo jest
 * stanem tak samo poprawnym jak wykaz niepusty.
 */
function wykazHostow(ksiazka: KsiazkaHostow, czynnosci: CzynnosciHosta): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-hosty';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-hosty__podpis';
  podpis.textContent = 'Książka hostów (wykaz widoku)';
  blok.append(podpis);

  const grupy = ksiazka.grupy();
  if (grupy.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Książka jest pusta. Wpisz adres celu i dołóż wpis albo wczytaj plik konfiguracyjny OpenSSH. ' +
      'Wpis dołożony przyciskiem idzie wprost do dziennika rdzenia i przeżywa odświeżenie strony ' +
      'oraz restart rdzenia; wpisy wczytane z pliku wnosi tam pozycja „Zapisz wykaz w rdzeniu”.';
    blok.append(pusto);
    return blok;
  }

  for (const [grupa, wpisy] of grupy) {
    const naglowek = document.createElement('h5');
    naglowek.className = 'dt-grupa';
    naglowek.textContent = grupa;
    blok.append(naglowek);

    const lista = wykaz(`Hosty folderu ${grupa}`, 'dt-wykaz');
    for (const wpis of wpisy) {
      const pozycja = pozycjaWykazu(wpis.nazwa, opisHosta(wpis), 'dt');

      const polacz = przyciskAkcji('Połącz', 'dn-btn dn-btn--atrament');
      polacz.title = `Otwiera kartę powłoki zdalnej do ${wpis.cel}. ${zdanieProgramuPowloki(TerminalShell.Ssh)}`;
      polacz.addEventListener('click', () => czynnosci.polacz(wpis));

      const wczytaj = przyciskAkcji('Wczytaj do formularza');
      wczytaj.addEventListener('click', () => czynnosci.wczytaj(wpis));

      const usun = przyciskAkcji('Usuń wpis');
      usun.addEventListener('click', () => czynnosci.usun(wpis));

      pozycja.akcje.append(polacz, wczytaj, usun);
      lista.append(pozycja.element);
    }
    blok.append(lista);
  }
  return blok;
}

/** Przekłada wpis książki hostów zwrócony przez rdzeń na odpowiadający mu wpis wykazu prowadzonego przez to okno. */
function wpisZKontraktu(host: TerminalHost): WpisHosta {
  const wpis: WpisHosta = {
    id: host.id,
    nazwa: host.name,
    cel: host.target,
    grupa: host.group ?? '',
    katalog: host.workingDir ?? '',
    notatka: host.note ?? '',
  };
  if (host.port !== undefined) wpis.port = host.port;
  return wpis;
}

/**
 * Wykaz kluczy SSH znanych rdzeniowi.
 *
 * Wiersz niesie ODCISK, a nie materiał klucza — i tak ma zostać: odcisk służy
 * porównaniu z tym, co pokazuje host docelowy, a klucz prywatny nie ma powodu
 * przechodzić przez łącze w żadną stronę.
 */
function wykazKluczy(
  klucze: readonly TerminalSshKey[],
  zdejmij: (klucz: TerminalSshKey) => void,
): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-klucze';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-klucze__podpis';
  podpis.textContent = 'Klucze SSH maszyny rdzenia';
  blok.append(podpis);

  if (klucze.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Wykaz kluczy jest pusty. Wytwórz klucz albo wciągnij do wykazu klucz leżący już na maszynie ' +
      'rdzenia — wpis książki hostów może potem wskazać go zamiast klucza domyślnego konfiguracji maszyny.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Klucze SSH', 'dt-wykaz');
  for (const klucz of klucze) {
    const czesci = [`rodzaj: ${klucz.keyType}`, `odcisk: ${klucz.fingerprint}`];
    czesci.push(klucz.hasPassphrase ? 'chroniony hasłem' : 'bez hasła');
    if (klucz.path !== undefined) czesci.push(`ścieżka: ${klucz.path}`);
    const pozycja = pozycjaWykazu(klucz.name, czesci.join(' · '), 'dt');

    if (klucz.publicKey !== undefined && klucz.publicKey !== '') {
      const wywoz = przyciskAkcji('Zapisz klucz publiczny');
      wywoz.title =
        'Zapisuje część JAWNĄ klucza — tę, którą przenosi się na host docelowy do pliku authorized_keys.';
      wywoz.addEventListener('click', () => {
        pobierzPlik(`${klucz.name}.pub`, `${klucz.publicKey ?? ''}\n`, 'text/plain');
      });
      pozycja.akcje.append(wywoz);
    }

    const zdejmijPrzycisk = przyciskAkcji('Zdejmij z wykazu');
    zdejmijPrzycisk.title =
      'Zdejmuje wpis z wykazu rdzenia. Pliki klucza na dysku maszyny rdzenia zostają nietknięte.';
    zdejmijPrzycisk.addEventListener('click', () => zdejmij(klucz));
    pozycja.akcje.append(zdejmijPrzycisk);
    lista.append(pozycja.element);
  }
  blok.append(lista);
  return blok;
}

/**
 * Wykaz przekierowań portów tunelu SSH wraz z ich bieżącym stanem oraz powodem niepowodzenia,
 * gdy założenie się nie udało.
 */
function wykazTuneli(
  tunele: readonly TerminalTunnel[],
  zamknij: (tunel: TerminalTunnel) => void,
): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-tunele';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-tunele__podpis';
  podpis.textContent = 'Przekierowania portów';
  blok.append(podpis);

  if (tunele.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Żadne przekierowanie portu nie zostało założone z tego okna. Tunel zakłada się polami adresu ' +
      'celu i portów w pasie narzędzi okna.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Tunele portowe', 'dt-wykaz');
  for (const tunel of tunele) {
    const czesci = [`rodzaj: ${tunel.kind}`, `stan: ${tunel.status}`];
    if (tunel.localPort !== undefined) czesci.push(`port maszyny rdzenia: ${tunel.localPort}`);
    if (tunel.remoteHost !== undefined && tunel.remoteHost !== '') {
      czesci.push(`cel po drugiej stronie: ${tunel.remoteHost}:${tunel.remotePort ?? '?'}`);
    }
    // Powód niepowodzenia wchodzi do opisu, bo sam stan failed nie mówi Operatorowi, co poprawić.
    if (tunel.errorMessage !== undefined && tunel.errorMessage !== '') {
      czesci.push(`powód: ${tunel.errorMessage}`);
    }
    const pozycja = pozycjaWykazu(tunel.id, czesci.join(' · '), 'dt');
    if (tunel.status === TerminalTunnelStatus.Active) {
      const zamknijPrzycisk = przyciskAkcji('Zamknij tunel');
      zamknijPrzycisk.addEventListener('click', () => zamknij(tunel));
      pozycja.akcje.append(zamknijPrzycisk);
    }
    lista.append(pozycja.element);
  }
  blok.append(lista);
  return blok;
}

/**
 * Opis wiersza hosta w wykazie: adres celu połączenia, katalog roboczy karty i notatka
 * pochodzenia wpisu.
 */
function opisHosta(wpis: WpisHosta): string {
  const czesci = [`cel: ${wpis.cel}`];
  czesci.push(`katalog: ${wpis.katalog === '' ? 'własny katalog okna' : wpis.katalog}`);
  if (wpis.notatka !== '') czesci.push(wpis.notatka);
  return czesci.join(' · ');
}

/**
 * Karty powłok znane rdzeniowi, uzupełnione o kartę otwartą w tym połączeniu, zanim zdąży wejść
 * do wykazu rdzenia; pierwszeństwo ma rdzeń, bo on rozstrzyga, czy karta wciąż istnieje.
 */
function wykazKart(zRdzenia: readonly TerminalSession[], stan: StanTerminala): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-sesje';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-sesje__podpis';
  podpis.textContent = 'Karty powłok prowadzone przez rdzeń';
  blok.append(podpis);

  const zastrzezenie = document.createElement('p');
  zastrzezenie.className = 'dn-pole-opis';
  zastrzezenie.textContent =
    'Wykaz pochodzi z rdzenia (terminal.session.list), więc niesie także karty otwarte przed ' +
    'rozłączeniem klienta — rdzeń odtwarza je przy starcie. Karty zamknięte nie wchodzą do tego wykazu.';
  blok.append(zastrzezenie);

  const zebrane = new Map<string, TerminalSession>();
  for (const karta of stan.karty()) zebrane.set(karta.id, karta);
  for (const karta of zRdzenia) zebrane.set(karta.id, karta);
  const karty = [...zebrane.values()];
  if (karty.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dn-pole-opis';
    pusto.textContent =
      'Rdzeń nie prowadzi ani jednej karty powłoki w tym oknie. Kartę otwiera się z wpisu książki ' +
      'hostów albo w oknie Terminal Tabs.';
    blok.append(pusto);
    return blok;
  }

  const lista = wykaz('Karty powłok okna', 'dt-wykaz');
  for (const karta of karty) {
    lista.append(pozycjaWykazu(karta.title ?? karta.shell, opisKartyOkna(karta), 'dt').element);
  }
  blok.append(lista);
  return blok;
}

function opisKartyOkna(karta: TerminalSession): string {
  const czesci = [
    `powłoka: ${karta.shell}`,
    `stan: ${karta.status === TerminalSessionStatus.Running ? 'karta czynna' : karta.status}`,
    `katalog: ${karta.workingDir ?? 'własny katalog okna'}`,
  ];
  if (karta.pid !== undefined) czesci.push(`PID powłoki: ${karta.pid}`);
  return czesci.join(' · ');
}

/**
 * Otwarcie karty powłoki zdalnej dla wskazanego wpisu książki hostów; źródło danych i stan
 * treści wchodzą parametrem wywołania.
 */
function polaczZHostem(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  tresc: StanTresci,
  wpis: WpisHosta,
): void {
  const okno = stan.okno();
  if (okno === '') {
    tresc.blad('Moduł nie zna okna komunikacji — karty zdalnej nie ma gdzie otworzyć.');
    return;
  }
  // Adres jedzie polem kontraktu, zmienna środowiska zostaje jako droga zastępcza dla starszego rdzenia.
  const zadanie: TerminalSessionOpenRequest =
    wpis.id !== undefined && wpis.id !== ''
      ? {
          windowId: okno,
          shell: TerminalShell.Ssh,
          title: wpis.nazwa === '' ? wpis.cel : wpis.nazwa,
          hostId: wpis.id,
        }
      : {
          windowId: okno,
          shell: TerminalShell.Ssh,
          title: wpis.nazwa === '' ? wpis.cel : wpis.nazwa,
          remoteTarget: wpis.cel,
          environment: [{ name: ZMIENNA_CELU, value: wpis.cel, enabled: true }],
          ...(wpis.port === undefined ? {} : { remotePort: wpis.port }),
          ...(wpis.katalog === '' ? {} : { workingDir: wpis.katalog }),
        };
  tresc.ladowanie(`Otwieranie karty powłoki zdalnej do ${wpis.cel}…`);
  void zrodlo.otworzKarte(zadanie).then((wynik) => {
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(
        `Rdzeń nie otworzył karty zdalnej do ${wpis.cel} (${Command.TerminalSessionOpen}). ` +
          zdanieProgramuPowloki(TerminalShell.Ssh),
        wynik.blad,
      );
      return;
    }
    const karta = wynik.wynik;
    stan.dodajKarte(karta);
    tresc.potwierdzenie(
      `Rdzeń otworzył kartę ${karta.id} do ${wpis.cel} — stan karty: ${karta.status}. ` +
        'Karta jest profilem połączenia: uwierzytelnienie i sam ruch po sieci rozstrzygną się dopiero ' +
        'przy pierwszym poleceniu.',
      karta.status === TerminalSessionStatus.Running,
    );
  });
}

/**
 * Zależność powłoki zdalnej dla noty okna.
 *
 * Wykaz bywa pusty i to też jest odpowiedź: gdyby rdzeń stracił program powłoki
 * zdalnej ze swojego wykazu wykonawczego, nota pokazałaby brak zamiast nazwy
 * programu, którego rdzeń nie uruchamia.
 */
function zaleznosciPowlokiZdalnej(): ProgramZewnetrzny[] {
  const zaleznosc = programPowloki(TerminalShell.Ssh);
  return zaleznosc === null ? [] : [zaleznosc];
}

/**
 * Kontrolki formularza i przycisków okna Session Manager: wpis hosta, klucze SSH, tunele
 * portowe oraz plik konfiguracyjny OpenSSH.
 */
interface PowierzchniaHostow {
  nazwa: HTMLInputElement;
  cel: HTMLInputElement;
  port: HTMLInputElement;
  grupa: HTMLInputElement;
  katalog: HTMLInputElement;
  notatka: HTMLInputElement;
  klucz: HTMLSelectElement;
  dodaj: HTMLButtonElement;
  polacz: HTMLButtonElement;
  eksport: HTMLButtonElement;
  plik: HTMLInputElement;
  zapiszWszystkie: HTMLButtonElement;
  nazwaKlucza: HTMLInputElement;
  rodzajKlucza: HTMLSelectElement;
  sciezkaKlucza: HTMLInputElement;
  wytworzKlucz: HTMLButtonElement;
  wciagnijKlucz: HTMLButtonElement;
  rodzajTunelu: HTMLSelectElement;
  hostDocelowy: HTMLInputElement;
  portZdalny: HTMLInputElement;
  portLokalny: HTMLInputElement;
  otworzTunel: HTMLButtonElement;
}

/**
 * Wpis hosta złożony z pól formularza wpisu; pusty adres celu połączenia oznacza brak wpisu
 * do zapisania.
 */
function wpisZFormularza(kontrolki: PowierzchniaHostow): WpisHosta | null {
  const cel = kontrolki.cel.value.trim();
  if (cel === '') return null;
  const nazwa = kontrolki.nazwa.value.trim();
  const wpis: WpisHosta = {
    nazwa: nazwa === '' ? cel : nazwa,
    cel,
    grupa: kontrolki.grupa.value.trim(),
    katalog: kontrolki.katalog.value.trim(),
    notatka: kontrolki.notatka.value.trim(),
  };
  const port = Number.parseInt(kontrolki.port.value, 10);
  // Port niepoprawny nie wchodzi wcale: rdzeń bierze wtedy port domyślny protokołu połączenia.
  if (Number.isFinite(port) && port > 0 && port <= 65535) wpis.port = port;
  return wpis;
}

/**
 * Przepisuje wpis książki do pól formularza — droga do poprawienia istniejącego wpisu bez
 * wpisywania go od nowa.
 */
function wpiszDoFormularza(kontrolki: PowierzchniaHostow, wpis: WpisHosta): void {
  kontrolki.nazwa.value = wpis.nazwa;
  kontrolki.cel.value = wpis.cel;
  kontrolki.port.value = wpis.port === undefined ? '' : String(wpis.port);
  kontrolki.grupa.value = wpis.grupa;
  kontrolki.katalog.value = wpis.katalog;
  kontrolki.notatka.value = wpis.notatka;
  kontrolki.cel.focus();
}

/**
 * Składa kontrolki, pasek akcji, pasek narzędzi i ciało okna; pozycje niewykonywane stoją
 * nieczynne z powodem liczonym z wykazu komend rdzenia.
 */
function zlozPowierzchnieHostow(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  pokrycie: PokrycieKomend,
): PowierzchniaHostow {
  const nazwa = pole('Nazwa wpisu hosta', 'np. web-01');
  const cel = pole('Adres celu połączenia', 'użytkownik@host albo alias konfiguracji serwera');
  const port = poleLiczbowe('Port połączenia', 'puste = port domyślny protokołu');
  // Wykaz kluczy wypełnia się po odczycie z rdzenia; pozycja pusta znaczy klucz domyślny maszyny.
  const klucz = wybor('Klucz SSH wpisu', [['', 'klucz domyślny konfiguracji maszyny rdzenia']]);
  const grupa = pole('Folder książki', 'np. Produkcja');
  const katalog = pole('Katalog roboczy karty zdalnej', 'puste = własny katalog okna');
  const notatka = pole('Notatka wpisu', 'np. host wdrożeniowy');

  const dodaj = przyciskAkcji('+ Dodaj host');
  const polacz = przyciskAkcji('Połącz — nowa karta SSH', 'dn-btn dn-btn--atrament');
  const eksport = przyciskAkcji('Eksportuj książkę hostów');

  // Pole wyboru pliku jest jedyną drogą, którą treść z maszyny Operatora wchodzi do przeglądarki.
  const plik = document.createElement('input');
  plik.type = 'file';
  plik.className = 'dn-pole-kontrolka dt-plik';
  plik.accept = '.config,.txt,text/plain';
  plik.setAttribute('aria-label', 'Wczytaj plik konfiguracyjny OpenSSH');

  const zapiszWszystkie = przyciskAkcji('Zapisz wykaz w rdzeniu');
  zapiszWszystkie.title =
    'Wnosi do książki rdzenia te wpisy wykazu, które nie mają jeszcze swojego wiersza — ' +
    'zwykle wpisy wczytane z pliku konfiguracyjnego Operatora.';

  const nazwaKlucza = pole('Nazwa klucza SSH', 'np. klucz-wydania');
  const rodzajKlucza = wybor('Rodzaj klucza SSH', [
    [TerminalKeyType.Ed25519, 'ed25519 — krótki i szybki, wybór domyślny'],
    [TerminalKeyType.Ecdsa, 'ecdsa — krzywa P-256'],
    [TerminalKeyType.Rsa, 'rsa — 3072 bity, gdy host nie zna nowszych'],
  ]);
  const sciezkaKlucza = pole(
    'Ścieżka klucza na maszynie rdzenia',
    'np. /home/operator/.ssh/id_ed25519 — do wciągnięcia do wykazu',
  );
  const wytworzKlucz = przyciskAkcji('Wytwórz klucz SSH');
  wytworzKlucz.title =
    'Wytwarza parę kluczy NA MASZYNIE RDZENIA. Część tajna nie wraca do przeglądarki ani teraz, ani później.';
  const wciagnijKlucz = przyciskAkcji('Wciągnij klucz do wykazu');
  wciagnijKlucz.title =
    'Wciąga do wykazu klucz leżący już na maszynie rdzenia. Klucz wskazuje się ŚCIEŻKĄ, nie treścią.';

  const rodzajTunelu = wybor('Rodzaj przekierowania portu', [
    [TerminalTunnelKind.Local, 'miejscowe (-L) — port rdzenia prowadzi do maszyny zdalnej'],
    [TerminalTunnelKind.Remote, 'zwrotne (-R) — port maszyny zdalnej prowadzi do rdzenia'],
    [TerminalTunnelKind.Dynamic, 'dynamiczne (-D) — pośrednik SOCKS'],
  ]);
  const hostDocelowy = pole('Maszyna docelowa tunelu', 'puste = localhost po drugiej stronie');
  const portZdalny = poleLiczbowe('Port docelowy tunelu', 'np. 5432');
  const portLokalny = poleLiczbowe('Port po stronie rdzenia', 'puste = rdzeń wybierze wolny');
  const otworzTunel = przyciskAkcji('Załóż tunel portowy');
  otworzTunel.title =
    'Zakłada przekierowanie programem ssh na maszynie rdzenia. Stan tunelu bierze się z tego procesu: ' +
    'przekierowanie, którego nie udało się założyć, wraca jako niepowodzenie wraz z powodem.';

  // Nazwa spoza kontraktu jest tu wskazaniem, nie zapisem stanu; odcisk nie wchodzi do scalenia.
  const znaneHosty = pokrycie.przycisk(
    'Polityka known_hosts',
    'terminal.knownhosts.get',
    'Podgląd odcisków kluczy hostów znanych maszynie rdzenia',
  );

  oznaczWarstwy([
    [polacz, 'zawsze'],
    [dodaj, 'na-zadanie'],
    [zapiszWszystkie, 'na-zadanie'],
    [eksport, 'kontekstowa'],
    [plik, 'kontekstowa'],
    [otworzTunel, 'ekspercka'],
    [wytworzKlucz, 'ekspercka'],
    [wciagnijKlucz, 'ekspercka'],
    [znaneHosty, 'ekspercka'],
    [nazwaKlucza, 'ekspercka'],
    [rodzajKlucza, 'ekspercka'],
    [sciezkaKlucza, 'ekspercka'],
    [rodzajTunelu, 'ekspercka'],
    [hostDocelowy, 'ekspercka'],
    [portZdalny, 'ekspercka'],
    [portLokalny, 'ekspercka'],
  ]);

  rama.akcje.append(
    polacz,
    dodaj,
    zapiszWszystkie,
    eksport,
    wytworzKlucz,
    wciagnijKlucz,
    otworzTunel,
    znaneHosty,
  );
  rama.narzedzia.append(
    nazwa,
    cel,
    port,
    grupa,
    katalog,
    notatka,
    klucz,
    plik,
    nazwaKlucza,
    rodzajKlucza,
    sciezkaKlucza,
    rodzajTunelu,
    hostDocelowy,
    portZdalny,
    portLokalny,
  );
  rama.cialo.append(
    notaZaleznosci(
      'Połączenie zdalne wykonuje program klienta OpenSSH leżący na maszynie rdzenia, nie przeglądarka ' +
        'i nie serwer aplikacji. Port i wskazanie klucza rdzeń czyta z wpisu książki i podaje temu ' +
        'programowi przełącznikami; wpis bez wskazań zostawia obie rzeczy konfiguracji jego maszyny.',
      zaleznosciPowlokiZdalnej(),
    ),
    stanTresci,
  );

  return {
    nazwa,
    cel,
    port,
    grupa,
    katalog,
    notatka,
    klucz,
    dodaj,
    polacz,
    eksport,
    plik,
    zapiszWszystkie,
    nazwaKlucza,
    rodzajKlucza,
    sciezkaKlucza,
    wytworzKlucz,
    wciagnijKlucz,
    rodzajTunelu,
    hostDocelowy,
    portZdalny,
    portLokalny,
    otworzTunel,
  };
}
