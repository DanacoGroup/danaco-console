/** Panel dopełnia Queue Managera dwunastoma czynnościami rodziny komend kolejki: zleceniami, polityką, zadaniami martwymi i głębokością w czasie. */
import { QueueAction, QueueItemStatus, QueueScope } from '../../../../shared/contract';
import { pole, poleLiczbowe, poleTresci, wybor } from '../../modele/kontrolki-formularza';
import {
  odczytajChwile,
  odczytajLiczbe,
  odczytajWykaz,
  odczytajZapis,
  utworzPanelDobudowy,
  wykonajCzynnoscPanelu,
} from './panel-dobudowy';
import type { StanAutomatyki } from './stan-automatyki';
import type { ZrodloAutomations } from './zrodlo-automations';

/** Interfejs opisuje panel dopełniający Queue Managera: element gotowy do osadzenia w oknie obok wykazu kolejek. */
export interface PanelZlecen {
  element: HTMLElement;
}

/** Stała wylicza stany zlecenia w wykazie zawężenia wraz z pozycją pustą, która znaczy wszystkie stany naraz. */
const STANY_ZLECENIA: ReadonlyArray<readonly [string, string]> = [
  ['', 'wszystkie stany'],
  ...Object.values(QueueItemStatus).map((stan): readonly [string, string] => [stan, stan]),
];

/** Stała wylicza działania silnika kolejek wprost z wyliczenia kontraktu, tworząc wykaz zupełny bez pominięć. */
const DZIALANIA_SILNIKA: ReadonlyArray<readonly [string, string]> = Object.values(
  QueueAction,
).map((czynnosc): readonly [string, string] => [czynnosc, czynnosc]);

/** Stała wylicza zasięgi kolejki wprost z wyliczenia kontraktu, gotowe do wskazania w polu wyboru panelu. */
const ZASIEGI: ReadonlyArray<readonly [string, string]> = Object.values(QueueScope).map(
  (zasieg): readonly [string, string] => [zasieg, zasieg],
);

export function utworzPanelZlecen(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
): PanelZlecen {
  const panel = utworzPanelDobudowy(
    'Zlecenia kolejki',
    'Jedenaście działań silnika kolejek na pojedynczym zleceniu wraz z wykazem ' +
      'zleceń, polityką kolejki, zadaniami martwymi i głębokością kolejki w czasie.',
  );

  const idKolejki = panel.dodajPole(
    'Kolejka',
    pole('Kolejka', 'pusta bierze kolejkę bieżącą modułu'),
  );
  const idZlecenia = panel.dodajPole('Zlecenie', pole('Zlecenie', 'zlec-…'));
  const zleceniaScalane = panel.dodajPole(
    'Zlecenia scalane',
    pole('Zlecenia scalane', 'zlec-1, zlec-2'),
  );
  const stanZlecenia = panel.dodajPole('Stan', wybor('Stan zlecenia', STANY_ZLECENIA));
  const zasieg = panel.dodajPole('Zasięg kolejki', wybor('Zasięg kolejki', ZASIEGI));
  const priorytet = panel.dodajPole('Priorytet', poleLiczbowe('Priorytet', '0'));
  const sekundy = panel.dodajPole('Sekundy', poleLiczbowe('Sekundy', '60'));
  const termin = panel.dodajPole('Termin', pole('Termin', '2026-08-20T07:00:00Z'));
  const kluczIdempotencji = panel.dodajPole(
    'Klucz idempotencji',
    pole('Klucz idempotencji', 'zapobiega podwójnemu wykonaniu'),
  );
  const kolejkaDocelowa = panel.dodajPole('Kolejka docelowa', pole('Kolejka docelowa', ''));
  const ekspertDocelowy = panel.dodajPole('Ekspert docelowy', pole('Ekspert docelowy', ''));
  const warunek = panel.dodajPole(
    'Warunek',
    pole('Warunek', 'pusty zdejmuje warunek zlecenia'),
  );
  const ladunek = panel.dodajPole(
    'Ładunek',
    poleTresci('Ładunek', 4, '{ }'),
    'Nośnik ładunku zlecenia, ładunków podziału, torów rozgałęzienia i polityki kolejki.',
  );

  /** Kolejka wskazana w polu, a gdy pole puste — kolejka bieżąca modułu. */
  function kolejka(): string | null {
    const wskazana = idKolejki.value.trim() === '' ? stan.kolejka() : idKolejki.value.trim();
    if (wskazana === '') {
      panel.tresc.potwierdzenie(
        'Nie ma jeszcze kolejki. Załóż ją w Queue Managerze albo wskaż jej identyfikator.',
        false,
      );
      return null;
    }
    return wskazana;
  }

  /** Zlecenie wskazane w polu; brak wstrzymuje czynność z nazwanym powodem. */
  function zlecenie(): string | null {
    if (idZlecenia.value.trim() === '') {
      panel.tresc.potwierdzenie(
        'Wskaż zlecenie — te czynności dotyczą pojedynczego zlecenia, nie kolejki jako całości.',
        false,
      );
      return null;
    }
    return idZlecenia.value.trim();
  }

  /** Zapis z pola ładunku; nieczytelny wstrzymuje czynność. */
  function zapisLadunku(): unknown | undefined | null {
    const odczytany = odczytajZapis(ladunek.value);
    if (odczytany === null) {
      panel.tresc.potwierdzenie('Ładunek w polu jest nieczytelny jako zapis strukturalny.', false);
    }
    return odczytany;
  }

  // Działanie na kolejce działa bez wskazania automatyki: obejmuje też kolejkę przeglądu ręcznego.
  const dzialanie = panel.dodajPole(
    'Działanie na kolejce',
    wybor('Działanie na kolejce', DZIALANIA_SILNIKA),
  );

  panel.dodajCzynnosc('Wykonaj na kolejce', () => {
    const kod = kolejka();
    if (kod === null) return;
    const zadanie: Parameters<ZrodloAutomations['dzialanieNaKolejce']>[0] = {
      queueId: kod, action: dzialanie.value as QueueAction,
    };
    if (idZlecenia.value.trim() !== '') zadanie.itemId = idZlecenia.value.trim();
    wykonajCzynnoscPanelu(panel, 'Działanie na kolejce…', zrodlo.dzialanieNaKolejce(zadanie),
      'Rdzeń odmówił działania na kolejce.',
      'Rdzeń wykonał działanie na kolejce; stan po działaniu niesie odpowiedź, ' +
        'nie etykieta przycisku.');
  });

  panel.dodajCzynnosc('Dołóż zlecenie', () => {
    const kod = kolejka();
    if (kod === null) return;
    const tresc = zapisLadunku();
    if (tresc === null) return;
    const zadanie: Parameters<ZrodloAutomations['dodajZlecenie']>[0] = {
      queueId: kod, payload: tresc ?? {},
    };
    const waga = odczytajLiczbe(priorytet.value);
    if (waga !== undefined) zadanie.priority = waga;
    const chwila = odczytajChwile(termin.value);
    if (chwila !== undefined) zadanie.scheduledAt = chwila;
    if (kluczIdempotencji.value.trim() !== '') {
      zadanie.idempotencyKey = kluczIdempotencji.value.trim();
    }
    wykonajCzynnoscPanelu(panel, 'Dokładanie zlecenia…', zrodlo.dodajZlecenie(zadanie),
      'Rdzeń nie dołożył zlecenia do kolejki.',
      'Rdzeń przyjął zlecenie; pole „duplicate” mówi, czy klucz idempotencji już był użyty.');
  });

  panel.dodajCzynnosc('Zdejmij zlecenie', () => {
    const kod = kolejka();
    const zlec = zlecenie();
    if (kod === null || zlec === null) return;
    wykonajCzynnoscPanelu(panel, 'Zdejmowanie zlecenia…',
      zrodlo.zdejmijZlecenie({ queueId: kod, itemId: zlec }),
      'Rdzeń nie zdjął zlecenia.',
      'Rdzeń odpowiedział na zdjęcie; pole „removed” mówi, czy zlecenie dało się jeszcze zdjąć.');
  });

  panel.dodajCzynnosc('Odłóż zlecenie', () => {
    const kod = kolejka();
    const zlec = zlecenie();
    if (kod === null || zlec === null) return;
    const odstep = odczytajLiczbe(sekundy.value);
    if (odstep === undefined) {
      panel.tresc.potwierdzenie('Odłożenie wymaga liczby sekund w polu „Sekundy”.', false);
      return;
    }
    wykonajCzynnoscPanelu(panel, 'Odkładanie zlecenia…',
      zrodlo.odlozZlecenie({ queueId: kod, itemId: zlec, delaySeconds: odstep }),
      'Rdzeń nie odłożył zlecenia.', 'Rdzeń odłożył wykonanie zlecenia.');
  });

  panel.dodajCzynnosc('Podziel zlecenie', () => {
    const kod = kolejka();
    const zlec = zlecenie();
    if (kod === null || zlec === null) return;
    const tresc = zapisLadunku();
    if (tresc === null || tresc === undefined) {
      panel.tresc.potwierdzenie(
        'Podział wymaga ładunków podzadań — wykazu w polu „Ładunek”, po jednym na zlecenie.',
        false,
      );
      return;
    }
    wykonajCzynnoscPanelu(panel, 'Podział zlecenia…',
      zrodlo.podzielZlecenie({ queueId: kod, itemId: zlec, payloads: tresc }),
      'Rdzeń nie podzielił zlecenia.',
      'Rdzeń założył podzadania; zlecenie źródłowe zostało zamknięte.');
  });

  panel.dodajCzynnosc('Scal zlecenia', () => {
    const kod = kolejka();
    if (kod === null) return;
    const scalane = odczytajWykaz(zleceniaScalane.value);
    if (scalane.length === 0) {
      panel.tresc.potwierdzenie('Wskaż zlecenia scalane w polu „Zlecenia scalane”.', false);
      return;
    }
    const tresc = zapisLadunku();
    if (tresc === null) return;
    const zadanie: Parameters<ZrodloAutomations['scalZlecenia']>[0] = {
      queueId: kod, itemIds: scalane,
    };
    if (tresc !== undefined) zadanie.payload = tresc;
    wykonajCzynnoscPanelu(panel, 'Scalanie zleceń…', zrodlo.scalZlecenia(zadanie),
      'Rdzeń nie scalił zleceń.',
      'Rdzeń założył zlecenie scalone; zlecenia źródłowe zostały zamknięte.');
  });

  panel.dodajCzynnosc('Skieruj zlecenie', () => {
    const kod = kolejka();
    const zlec = zlecenie();
    if (kod === null || zlec === null) return;
    const zadanie: Parameters<ZrodloAutomations['skierujZlecenie']>[0] = {
      queueId: kod, itemId: zlec,
    };
    if (kolejkaDocelowa.value.trim() !== '') zadanie.targetQueueId = kolejkaDocelowa.value.trim();
    if (ekspertDocelowy.value.trim() !== '') zadanie.targetAgentId = ekspertDocelowy.value.trim();
    wykonajCzynnoscPanelu(panel, 'Kierowanie zlecenia…', zrodlo.skierujZlecenie(zadanie),
      'Rdzeń nie skierował zlecenia.',
      'Rdzeń skierował zlecenie do wskazanej kolejki albo do wskazanego eksperta.');
  });

  panel.dodajCzynnosc('Rozgałęź zlecenie', () => {
    const kod = kolejka();
    const zlec = zlecenie();
    if (kod === null || zlec === null) return;
    const tresc = zapisLadunku();
    if (tresc === null || tresc === undefined) {
      panel.tresc.potwierdzenie(
        'Rozgałęzienie wymaga torów — wykazu w polu „Ładunek”, po jednym zapisie na tor.',
        false,
      );
      return;
    }
    const tory = tresc as Parameters<ZrodloAutomations['rozgalezZlecenie']>[0]['branches'];
    wykonajCzynnoscPanelu(panel, 'Rozgałęzianie zlecenia…',
      zrodlo.rozgalezZlecenie({ queueId: kod, itemId: zlec, branches: tory }),
      'Rdzeń nie rozgałęził zlecenia.', 'Rdzeń założył zlecenia torów równoległych.');
  });

  panel.dodajCzynnosc('Uwarunkuj zlecenie', () => {
    const kod = kolejka();
    const zlec = zlecenie();
    if (kod === null || zlec === null) return;
    wykonajCzynnoscPanelu(panel, 'Zapis warunku…',
      zrodlo.uwarunkujZlecenie({ queueId: kod, itemId: zlec, condition: warunek.value.trim() }),
      'Rdzeń nie zapisał warunku zlecenia.',
      'Rdzeń zapisał warunek przetworzenia; warunek pusty go zdejmuje.');
  });

  panel.dodajCzynnosc('Wykaz zleceń', () => {
    const kod = kolejka();
    if (kod === null) return;
    const zadanie: Parameters<ZrodloAutomations['zleceniaKolejki']>[0] = { queueId: kod };
    if (stanZlecenia.value !== '') zadanie.status = stanZlecenia.value as QueueItemStatus;
    wykonajCzynnoscPanelu(panel, 'Odczyt zleceń…', zrodlo.zleceniaKolejki(zadanie),
      'Rdzeń nie oddał zleceń kolejki.',
      'Rdzeń oddał zlecenia w porządku przetwarzania wraz z liczbą wszystkich.');
  });

  panel.dodajCzynnosc('Zapisz politykę', () => {
    const kod = kolejka();
    if (kod === null) return;
    const tresc = zapisLadunku();
    if (tresc === null) return;
    const polityka = {
      ...((tresc ?? {}) as Record<string, unknown>),
      scope: zasieg.value as QueueScope,
    } as Parameters<ZrodloAutomations['ustawPolitykeKolejki']>[0]['policy'];
    wykonajCzynnoscPanelu(panel, 'Zapis polityki kolejki…',
      zrodlo.ustawPolitykeKolejki({ queueId: kod, policy: polityka }),
      'Rdzeń nie zapisał polityki kolejki.',
      'Rdzeń zapisał zasięg, współbieżność, przepustowość i politykę ponawiania.');
  });

  panel.dodajCzynnosc('Zadania martwe', () => {
    const zadanie: Parameters<ZrodloAutomations['zadaniaMartwe']>[0] = {};
    if (idKolejki.value.trim() !== '') zadanie.queueId = idKolejki.value.trim();
    wykonajCzynnoscPanelu(panel, 'Odczyt zadań martwych…', zrodlo.zadaniaMartwe(zadanie),
      'Rdzeń nie oddał zadań martwych.',
      'Rdzeń oddał zlecenia trwale nieudane; bez wskazania kolejki — ze wszystkich kolejek.');
  });

  panel.dodajCzynnosc('Głębokość kolejki', () => {
    const odcinek = odczytajLiczbe(sekundy.value) ?? 3600;
    const od = odczytajChwile(termin.value) ?? Date.now() - 7 * 24 * 3600 * 1000;
    const zadanie: Parameters<ZrodloAutomations['glebokoscKolejki']>[0] = {
      fromAt: od, bucketSeconds: odcinek,
    };
    if (idKolejki.value.trim() !== '') zadanie.queueId = idKolejki.value.trim();
    else if (zasieg.value !== '') zadanie.scope = zasieg.value as QueueScope;
    wykonajCzynnoscPanelu(panel, 'Rachunek głębokości…', zrodlo.glebokoscKolejki(zadanie),
      'Rdzeń nie oddał głębokości kolejki.',
      'Rdzeń oddał liczbę zleceń oczekujących w kolejnych odcinkach czasu.');
  });

  return { element: panel.element };
}
