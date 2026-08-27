import type { LibraryFile } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { StanBiblioteki } from './stan-biblioteki';
import type { ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Ścieżki zapisu czynności zbiorczych, osobno od odczytu i widoku: każda
 * sprawdza warunek wstępny, wysyła komendę i oddaje jedno zdanie odpowiedzi.
 */
export interface OdpowiedzZapisu {
  tresc: string;
  powodzenie: boolean;
}

const BEZ_ZAZNACZENIA = 'Zaznacz plik w wykazie — czynność zbiorcza działa na zaznaczeniu.';

/**
 * Etykieta zbiorcza wysyła komplet etykiet pliku po zmianie, sumując
 * dotychczasowe z nową, a potwierdzenie bierze etykiety z odpowiedzi rdzenia,
 * nie z żądania, żeby ubytek był ogłaszany, nie przemilczany.
 */
export async function nadajEtykieteZbiorczo(
  stan: StanBiblioteki,
  etykieta: string,
  pliki: readonly LibraryFile[],
): Promise<OdpowiedzZapisu> {
  if (pliki.length === 0) return { tresc: BEZ_ZAZNACZENIA, powodzenie: false };
  const nowa = etykieta.trim();
  if (nowa === '') {
    return {
      tresc:
        'Pole etykiety jest puste po przycięciu przez okno — nic nie zostało wysłane. ' +
        'Wpisz etykietę.',
      powodzenie: false,
    };
  }
  const przyciete = nowa !== etykieta;
  const bezEtykiety: string[] = [];
  const zUbytkiem: string[] = [];

  for (const plik of pliki) {
    const przed = plik.tags ?? [];
    const komplet = [...new Set([...przed, nowa])];
    const wynik = await stan.zrodlo.nadajEtykiety(plik.id, komplet);
    if (!wynik.udany || wynik.wynik === undefined) {
      return {
        tresc: opisOdmowy(`Etykieta pliku „${plik.name}"`, wynik.blad?.code, wynik.blad?.message),
        powodzenie: false,
      };
    }
    const po = wynik.wynik.file.tags ?? [];
    stan.wchlon(wynik.wynik.file);
    if (!po.includes(nowa)) bezEtykiety.push(plik.name);
    const utracone = przed.filter((zastana) => !po.includes(zastana));
    if (utracone.length > 0) zUbytkiem.push(`${plik.name} (${utracone.join(', ')})`);
  }

  const dopisek = przyciete ? ` Okno wysłało „${nowa}" zamiast wpisanego „${etykieta}".` : '';
  if (bezEtykiety.length > 0) {
    return {
      tresc:
        `Rdzeń przyjął komendę, ale po zapisie etykiety „${nowa}" nie ma w plikach: ` +
        `${bezEtykiety.join(', ')}. Nadania nie potwierdzam.${dopisek}`,
      powodzenie: false,
    };
  }
  if (zUbytkiem.length > 0) {
    return {
      tresc:
        `Etykieta „${nowa}" stoi po zapisie przy plikach: ${pliki.length}, ale rdzeń oddał ` +
        `je bez etykiet zastanych — ${zUbytkiem.join('; ')}. Zapis podmienia komplet, ` +
        `więc odśwież wykaz przed kolejnym nadaniem.${dopisek}`,
      powodzenie: false,
    };
  }
  return {
    tresc: `Etykieta „${nowa}" stoi po zapisie przy plikach: ${pliki.length}.${dopisek}`,
    powodzenie: true,
  };
}

/**
 * Przypisanie zaznaczenia do kolekcji porównuje liczbę przypisanych z liczbą
 * zamówioną, bo rdzeń pomija pliki, których nie zna, a żądanie z nieznanymi
 * kodami mimo to wraca powodzeniem.
 */
export async function przypiszZaznaczenie(
  stan: StanBiblioteki,
  idKolekcji: string,
  pliki: readonly LibraryFile[],
): Promise<OdpowiedzZapisu> {
  if (pliki.length === 0) return { tresc: BEZ_ZAZNACZENIA, powodzenie: false };
  if (idKolekcji === '') {
    return {
      tresc: 'Brak kolekcji do wskazania — załóż ją w oknie Tags & Collections.',
      powodzenie: false,
    };
  }
  const wynik = await stan.zrodlo.przypiszDoKolekcji(
    idKolekcji,
    pliki.map((plik) => plik.id),
  );
  if (!wynik.udany || wynik.wynik === undefined) {
    return {
      tresc: opisOdmowy('Przypisanie do kolekcji', wynik.blad?.code, wynik.blad?.message),
      powodzenie: false,
    };
  }
  const przypisane = wynik.wynik.assignedCount;
  if (przypisane !== pliki.length) {
    return {
      tresc:
        `Do kolekcji ${idKolekcji} zamówiono zasobów: ${pliki.length}, a rdzeń przypisał: ` +
        `${przypisane}. Rozbieżności nie potwierdzam — odśwież wykaz i sprawdź, które ` +
        'pliki rdzeń jeszcze zna.',
      powodzenie: false,
    };
  }
  return {
    tresc: `Do kolekcji ${idKolekcji} przypisano zasobów: ${przypisane}.`,
    powodzenie: true,
  };
}

/**
 * Otwarcie zasobów w module docelowym bierze okno źródłowe z rejestru okien
 * komunikacji i odczytuje moduł docelowy z okna oddanego w odpowiedzi, nie
 * z pola formularza, bo tylko tam stoi moduł faktycznie nadany.
 */
export async function przeniesZasoby(
  stan: StanBiblioteki,
  otoczenie: ZrodloOtoczenia,
  modulWpisany: string,
  pliki: readonly LibraryFile[],
): Promise<OdpowiedzZapisu> {
  if (pliki.length === 0) return { tresc: BEZ_ZAZNACZENIA, powodzenie: false };
  if (stan.idOkna() === '') {
    return {
      tresc:
        'Przeniesienie kontekstu żąda okna źródłowego, a rdzeń nie oddał ani jednego okna ' +
        'komunikacji tej sesji (window.list).',
      powodzenie: false,
    };
  }
  const cel = modulWpisany.trim() || (pliki[0]?.sourceModuleId ?? '');
  if (cel === '') {
    return {
      tresc: 'Rdzeń nie podał modułu wytwórcy pliku — wskaż moduł docelowy wprost.',
      powodzenie: false,
    };
  }
  const wynik = await otoczenie.przeniesKontekst(stan.idOkna(), cel, {
    documentIds: pliki.map((plik) => plik.id),
  });
  if (!wynik.udany || wynik.wynik?.transferred !== true) {
    return {
      tresc: opisOdmowy('Otwarcie w module docelowym', wynik.blad?.code, wynik.blad?.message),
      powodzenie: false,
    };
  }
  const oknoDocelowe = wynik.wynik.window;
  if (oknoDocelowe === undefined) {
    return {
      tresc:
        `Rdzeń potwierdził przeniesienie zasobów: ${pliki.length}, ale nie oddał okna ` +
        'docelowego, więc nie wiem, dokąd komplet trafił.',
      powodzenie: false,
    };
  }
  if (oknoDocelowe.moduleId !== cel) {
    return {
      tresc:
        `Przeniesienie zamówiono do modułu ${cel}, a rdzeń oddał okno ${oknoDocelowe.id} ` +
        `modułu „${oknoDocelowe.moduleId}". Rozbieżności nie potwierdzam.`,
      powodzenie: false,
    };
  }
  return {
    tresc:
      `Zasoby (${pliki.length}) przeniesione — rdzeń wskazał okno ${oknoDocelowe.id} ` +
      `modułu ${oknoDocelowe.moduleId}.`,
    powodzenie: true,
  };
}
