/**
 * Panel Workflow Buildera obsługuje trzynaście czynności paneli wersji i zmiennych, stojąc
 * wewnątrz budowniczego, bo cała ta praca dotyczy definicji automatyki.
 */
import {
  pole,
  poleLiczbowe,
  poleTresci,
} from '../../modele/kontrolki-formularza';
import {
  odczytajLiczbe,
  odczytajWykaz,
  odczytajZapis,
  utworzPanelDobudowy,
  wykonajCzynnoscPanelu,
} from './panel-dobudowy';
import type { StanAutomatyki } from './stan-automatyki';
import type { ZrodloAutomations } from './zrodlo-automations';

/** Panel dopełniający Workflow Buildera, osadzany wewnątrz kreatora wraz z jego szufladami wersji i zmiennych. */
export interface PanelWersji {
  element: HTMLElement;
}

export function utworzPanelWersji(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
): PanelWersji {
  const panel = utworzPanelDobudowy(
    'Wersje, zmienne i szablony',
    'Panele Workflow Buildera: przebieg próbny, historia wersji, etykiety, ' +
      'zmienne przepływu, notatka i układ kanwy, biblioteka szablonów, publikacja ' +
      'i udostępnienie automatyki.',
  );

  const wersjaOd = panel.dodajPole('Wersja wyjściowa', poleLiczbowe('Wersja wyjściowa', '1'));
  const wersjaDo = panel.dodajPole('Wersja porównywana', poleLiczbowe('Wersja porównywana', '2'));
  const krok = panel.dodajPole('Krok', pole('Krok', 'krok-1'));
  const etykiety = panel.dodajPole(
    'Etykiety',
    pole('Etykiety', 'raport, tygodniowy'),
    'Komplet etykiet po zmianie; wykaz pusty zdejmuje wszystkie.',
  );
  const nazwa = panel.dodajPole('Nazwa szablonu albo automatyki', pole('Nazwa', 'Raport tygodniowy'));
  const szablon = panel.dodajPole('Szablon', pole('Szablon', 'szablon-…'));
  const zapis = panel.dodajPole(
    'Zapis strukturalny',
    poleTresci('Zapis strukturalny', 4, '{ }'),
    'Nośnik danych próbnych symulacji, zmiennych i mapowań, układu kanwy oraz ' +
      'wartości parametrów szablonu — kontrakt niesie je zapisem strukturalnym.',
  );

  /** Automatyka bieżąca modułu; pusta wstrzymuje czynność z nazwanym powodem. */
  function automatyka(): string | null {
    if (stan.automatyka() === '') {
      panel.tresc.potwierdzenie(
        'Wskaż automatykę w Workflow Builderze — te czynności dotyczą jej definicji.',
        false,
      );
      return null;
    }
    return stan.automatyka();
  }

  /** Zapis z pola; nieczytelny wstrzymuje czynność zamiast lecieć do rdzenia. */
  function zapisPola(): unknown | undefined | null {
    const odczytany = odczytajZapis(zapis.value);
    if (odczytany === null) {
      panel.tresc.potwierdzenie(
        'Zapis strukturalny w polu jest nieczytelny — popraw go, zanim żądanie pójdzie do rdzenia.',
        false,
      );
    }
    return odczytany;
  }

  panel.dodajCzynnosc('Przebieg próbny', () => {
    const kod = automatyka();
    if (kod === null) return;
    const probne = zapisPola();
    if (probne === null) return;
    const zadanie: Parameters<ZrodloAutomations['symulujPrzeplyw']>[0] = { workflowId: kod };
    if (probne !== undefined) zadanie.sampleInput = probne;
    if (krok.value.trim() !== '') zadanie.stopAtStepId = krok.value.trim();
    wykonajCzynnoscPanelu(panel, 'Przebieg próbny…', zrodlo.symulujPrzeplyw(zadanie),
      'Rdzeń odmówił przebiegu próbnego.',
      'Rdzeń przeprowadził przebieg próbny; efekty uboczne kroków były wstrzymane.');
  });

  panel.dodajCzynnosc('Wykaz wersji', () => {
    const kod = automatyka();
    if (kod === null) return;
    wykonajCzynnoscPanelu(panel, 'Odczyt wersji…',
      zrodlo.wersjeAutomatyki({ workflowId: kod }),
      'Rdzeń nie oddał wersji definicji.', 'Rdzeń oddał historię wersji definicji.');
  });

  panel.dodajCzynnosc('Przywróć wersję', () => {
    const kod = automatyka();
    if (kod === null) return;
    const numer = odczytajLiczbe(wersjaOd.value);
    if (numer === undefined) {
      panel.tresc.potwierdzenie('Wskaż numer wersji przywracanej w polu „Wersja wyjściowa”.', false);
      return;
    }
    wykonajCzynnoscPanelu(panel, 'Przywracanie wersji…',
      zrodlo.przywrocWersje({ workflowId: kod, version: numer }),
      'Rdzeń nie przywrócił wersji.',
      `Rdzeń przywrócił wersję ${numer}; wersja zastana została w historii.`);
  });

  panel.dodajCzynnosc('Porównaj wersje', () => {
    const kod = automatyka();
    if (kod === null) return;
    const od = odczytajLiczbe(wersjaOd.value);
    const doWersji = odczytajLiczbe(wersjaDo.value);
    if (od === undefined || doWersji === undefined) {
      panel.tresc.potwierdzenie('Porównanie wymaga obu numerów wersji.', false);
      return;
    }
    wykonajCzynnoscPanelu(panel, 'Porównanie wersji…',
      zrodlo.porownajWersje({ workflowId: kod, fromVersion: od, toVersion: doWersji }),
      'Rdzeń nie oddał różnicy wersji.', 'Rdzeń oddał różnicę strukturalną wersji.');
  });

  panel.dodajCzynnosc('Zapisz etykiety', () => {
    const kod = automatyka();
    if (kod === null) return;
    wykonajCzynnoscPanelu(panel, 'Zapis etykiet…',
      zrodlo.ustawEtykiety({ workflowId: kod, tags: odczytajWykaz(etykiety.value) }),
      'Rdzeń nie zapisał etykiet.', 'Rdzeń zapisał komplet etykiet automatyki.');
  });

  panel.dodajCzynnosc('Zapisz zmienne i mapowania', () => {
    const kod = automatyka();
    if (kod === null) return;
    const odczytany = zapisPola();
    if (odczytany === null) return;
    const tresc = (odczytany ?? {}) as {
      variables?: Parameters<ZrodloAutomations['ustawZmienne']>[0]['variables'];
      mappings?: Parameters<ZrodloAutomations['ustawZmienne']>[0]['mappings'];
    };
    const zadanie: Parameters<ZrodloAutomations['ustawZmienne']>[0] = {
      workflowId: kod, variables: tresc.variables ?? [],
    };
    if (tresc.mappings !== undefined) zadanie.mappings = tresc.mappings;
    wykonajCzynnoscPanelu(panel, 'Zapis zmiennych…', zrodlo.ustawZmienne(zadanie),
      'Rdzeń nie zapisał zmiennych przepływu.',
      'Rdzeń zapisał zmienne przepływu; zastrzeżenia, o ile są, stoją w odpowiedzi.');
  });

  panel.dodajCzynnosc('Zapisz notatkę kroku', () => {
    const kod = automatyka();
    if (kod === null) return;
    if (krok.value.trim() === '') {
      panel.tresc.potwierdzenie('Wskaż krok — notatka należy do pojedynczego węzła.', false);
      return;
    }
    wykonajCzynnoscPanelu(panel, 'Zapis notatki…',
      zrodlo.ustawNotatkeKroku({
        workflowId: kod, stepId: krok.value.trim(), note: nazwa.value,
      }),
      'Rdzeń nie zapisał notatki kroku.',
      'Rdzeń zapisał notatkę kroku; treść pusta zdejmuje ją.');
  });

  panel.dodajCzynnosc('Zapisz układ kanwy', () => {
    const kod = automatyka();
    if (kod === null) return;
    const odczytany = zapisPola();
    if (odczytany === null) return;
    const polozenia = (odczytany ?? []) as Parameters<
      ZrodloAutomations['ustawUkladKanwy']
    >[0]['positions'];
    wykonajCzynnoscPanelu(panel, 'Zapis układu kanwy…',
      zrodlo.ustawUkladKanwy({ workflowId: kod, positions: polozenia }),
      'Rdzeń nie zapisał układu kanwy.',
      'Rdzeń zapisał położenia węzłów; układ przeżyje zamknięcie karty.');
  });

  panel.dodajCzynnosc('Zapisz jako szablon', () => {
    const kod = automatyka();
    if (kod === null) return;
    wykonajCzynnoscPanelu(panel, 'Zapis szablonu…',
      zrodlo.zapiszSzablon({ workflowId: kod, name: nazwa.value.trim() }),
      'Rdzeń nie zapisał szablonu przepływu.', 'Rdzeń zapisał definicję jako szablon.');
  });

  panel.dodajCzynnosc('Biblioteka szablonów', () => {
    wykonajCzynnoscPanelu(panel, 'Odczyt biblioteki…', zrodlo.szablony({}),
      'Rdzeń nie oddał biblioteki szablonów.', 'Rdzeń oddał bibliotekę szablonów przepływów.');
  });

  panel.dodajCzynnosc('Załóż z szablonu', () => {
    if (szablon.value.trim() === '') {
      panel.tresc.potwierdzenie('Wskaż szablon, z którego ma powstać automatyka.', false);
      return;
    }
    const wartosci = zapisPola();
    if (wartosci === null) return;
    const zadanie: Parameters<ZrodloAutomations['zastosujSzablon']>[0] = {
      templateId: szablon.value.trim(), name: nazwa.value.trim(),
    };
    if (wartosci !== undefined) zadanie.values = wartosci;
    wykonajCzynnoscPanelu(panel, 'Zakładanie z szablonu…', zrodlo.zastosujSzablon(zadanie),
      'Rdzeń nie założył automatyki z szablonu.',
      'Rdzeń założył automatykę z szablonu; parametry bez wartości nazywa odpowiedź.');
  });

  panel.dodajCzynnosc('Opublikuj wersję', () => {
    const kod = automatyka();
    if (kod === null) return;
    const zadanie: Parameters<ZrodloAutomations['opublikujAutomatyke']>[0] = { workflowId: kod };
    const numer = odczytajLiczbe(wersjaOd.value);
    if (numer !== undefined) zadanie.version = numer;
    wykonajCzynnoscPanelu(panel, 'Publikacja wersji…', zrodlo.opublikujAutomatyke(zadanie),
      'Rdzeń nie opublikował wersji.',
      'Rdzeń opublikował wersję; produkcyjnie wykonuje się odtąd ona, nie robocza.');
  });

  panel.dodajCzynnosc('Przełącz udostępnienie', () => {
    const kod = automatyka();
    if (kod === null) return;
    udostepniona = !udostepniona;
    wykonajCzynnoscPanelu(panel, 'Zmiana udostępnienia…',
      zrodlo.udostepnijAutomatyke({ workflowId: kod, shared: udostepniona }),
      'Rdzeń nie zmienił udostępnienia automatyki.',
      udostepniona
        ? 'Rdzeń przyjął udostępnienie automatyki w organizacji.'
        : 'Rdzeń przyjął zdjęcie udostępnienia automatyki.');
  });

  // Stan przełącznika trzymany w panelu, bo zapis przyjmuje wartość docelową, nie sam przełącz.
  let udostepniona = false;

  return { element: panel.element };
}
