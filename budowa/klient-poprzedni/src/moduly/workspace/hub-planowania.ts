import {
  WorkspaceCalendarSpan,
  WorkspaceProjectStatus,
  WorkspaceTaskStatus,
  type WorkspaceTask,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleTresci,
  pozycjaWykazu,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
  wykaz,
} from '../../modele/kontrolki-formularza';
import {
  kartyKolumny,
  postepProjektu,
  type CzynnosciPlanowania,
} from './czynnosci-planowania';
import type { CzynnosciWiedzy } from './czynnosci-wiedzy';
import type { StanProjektu } from './stan-projektu';
import { utworzStanTresci } from './stany-okna';

/** Hub planowania jest oknem zadań projektu w czterech widokach współdzielących jeden zbiór danych: liście, tablicy kanban, osi czasu i kalendarzu. */
export interface OknoHubuPlanowania {
  element: HTMLElement;
  odswiez(): void;
}

/** Widoki huba w kolejności zakładek opracowania obejmują listę, tablicę kanban, oś czasu Gantt oraz kalendarz zadań. */
const WIDOKI: ReadonlyArray<readonly [string, string]> = [
  ['lista', 'widok: Lista'],
  ['tablica', 'widok: Tablica kanban'],
  ['os', 'widok: Oś czasu'],
  ['kalendarz', 'widok: Kalendarz'],
];

/** Stany zadania pochodzą wprost z kontraktu — wykaz w widoku nie jest przepisany ręcznie, tylko wyliczony z ich wartości. */
const STANY_ZADANIA: ReadonlyArray<readonly [string, string]> = Object.values(
  WorkspaceTaskStatus,
).map((stan) => [stan, `stan: ${stan}`]);

/** Stany projektu pochodzą wprost z kontraktu, a plakietka panelu projektu przyjmuje dokładnie trzy takie wartości. */
const STANY_PROJEKTU: ReadonlyArray<readonly [string, string]> = Object.values(
  WorkspaceProjectStatus,
).map((stan) => [stan, `projekt: ${stan}`]);

export function utworzOknoHubuPlanowania(
  planowanie: CzynnosciPlanowania,
  wiedza: CzynnosciWiedzy,
  stan: StanProjektu,
): OknoHubuPlanowania {
  const rama = utworzRameOkna({
    tytul: 'Hub planowania',
    kod: 'planning-hub',
    rola: 'wiodące',
    przeznaczenie:
      'Zadania projektu w czterech widokach: lista, tablica kanban, oś czasu i kalendarz.',
    przedrostek: 'dw',
  });
  const tresc = utworzStanTresci();

  const tytul = pole('Tytuł zadania', 'np. Przygotować projekt pisma');
  const wskazaneZadanie = pole('Zadanie wskazane', 'identyfikator zadania z wykazu');
  const drugieZadanie = pole('Zadanie następujące', 'identyfikator zadania zależnego');
  const wskazanaZaleznosc = pole('Zależność', 'identyfikator zależności');
  const stanZadania = wybor('Stan zadania', STANY_ZADANIA);
  const stanProjektu = wybor('Stan projektu', STANY_PROJEKTU);
  const trescIcal = poleTresci('Treść pliku iCal', 3);
  const widok = wybor('Widok huba', WIDOKI);

  const zaloz = przycisk('+ Zadanie', 'dn-btn dn-btn--atrament');
  const zmien = przycisk('Zmień stan zadania');
  const przenies = przycisk('Przenieś kartę');
  const usun = przycisk('Usuń zadanie');
  const zwiaz = przycisk('Zwiąż zadania');
  const rozwiaz = przycisk('Znieś zależność');
  const wciagnij = przycisk('Wciągnij iCal');
  const zapiszIcal = przycisk('Zapisz iCal');
  const przestawProjekt = przycisk('Ustaw stan projektu');
  const odswiezOkno = przycisk('Odśwież widok');

  rama.akcje.append(
    zaloz,
    zmien,
    przenies,
    usun,
    zwiaz,
    rozwiaz,
    wciagnij,
    zapiszIcal,
    przestawProjekt,
    odswiezOkno,
  );
  rama.narzedzia.append(widok, stanZadania, stanProjektu);
  rama.cialo.append(
    wiersz('Tytuł nowego zadania', tytul, {
      klasa: 'dw-wiersz',
      objasnienie: 'Zadanie powstaje w stanie „todo”, w kolumnie tego stanu na tablicy.',
    }),
    wiersz('Zadanie wskazane', wskazaneZadanie, {
      klasa: 'dw-wiersz',
      objasnienie: 'Kliknięcie pozycji wykazu wpisuje tu jej identyfikator.',
    }),
    wiersz('Zadanie następujące', drugieZadanie, {
      klasa: 'dw-wiersz',
      objasnienie: 'Zależność „koniec–początek”: wskazane poprzedza następujące.',
    }),
    wiersz('Zależność', wskazanaZaleznosc, {
      klasa: 'dw-wiersz',
      objasnienie: 'Identyfikator z odpowiedzi po związaniu zadań — po nim znosi się zależność.',
    }),
    wiersz('Plik iCal', trescIcal, {
      klasa: 'dw-wiersz',
      objasnienie: 'Wydarzenia wchodzą jako zadania z terminem; pominięcia wracają z powodem.',
    }),
    tresc.element,
  );

  function projektAlboOstrzez(): string {
    const projekt = stan.projekt();
    if (projekt === '') {
      tresc.blad('Hub bez projektu nie ma zakresu — wskaż projekt na pulpicie.');
      return '';
    }
    return projekt;
  }

  function pozycjaZadania(zadanie: WorkspaceTask): HTMLElement {
    const opis = [
      zadanie.status,
      zadanie.assigneeId ?? 'bez wykonawcy',
      zadanie.dueAt === undefined ? 'bez terminu' : new Date(zadanie.dueAt).toISOString(),
    ].join(' · ');
    const { element, akcje } = pozycjaWykazu(zadanie.title, opis, 'dw');
    const wskaz = przycisk('Wskaż', 'dn-btn dn-btn--zarys');
    wskaz.addEventListener('click', () => {
      wskazaneZadanie.value = zadanie.id;
      tresc.potwierdzenie(`Wskazano zadanie „${zadanie.title}”.`, true);
    });
    akcje.append(wskaz);
    return element;
  }

  function pokazZadania(zadania: readonly WorkspaceTask[], naglowek: string): void {
    rama.ustawZnacznik(`postęp ${postepProjektu(zadania)}%`);
    if (zadania.length === 0) {
      tresc.pusto(`${naglowek}: projekt nie ma jeszcze zadań.`);
      return;
    }
    const lista = wykaz(naglowek, 'dw-wykaz');
    for (const zadanie of zadania) lista.append(pozycjaZadania(zadanie));
    tresc.tresc().append(lista);
  }

  function odczytajListe(projekt: string): void {
    void planowanie
      .zadania({ projectId: projekt, includeDone: true })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się odczytać zadań projektu.', wynik.blad);
          return;
        }
        pokazZadania(wynik.wynik, 'Zadania projektu');
      });
  }

  function odczytajTablice(projekt: string): void {
    void planowanie.tablica({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się odczytać tablicy projektu.', wynik.blad);
        return;
      }
      const tablica = wynik.wynik;
      rama.ustawZnacznik(`postęp ${postepProjektu(tablica.tasks)}%`);
      const miejsce = tresc.tresc();
      for (const kolumna of tablica.columns) {
        const karty = kartyKolumny(tablica, kolumna.id);
        const granica = kolumna.wipLimit ?? 0;
        const podpis =
          granica > 0 && karty.length > granica
            ? `${kolumna.name} (${karty.length}, granica ${granica} przekroczona)`
            : `${kolumna.name} (${karty.length})`;
        const lista = wykaz(podpis, 'dw-wykaz');
        lista.append(pozycjaWykazu(podpis, `stan ${kolumna.status}`, 'dw').element);
        for (const karta of karty) lista.append(pozycjaZadania(karta));
        miejsce.append(lista);
      }
      if (tablica.columns.length === 0) {
        tresc.pusto('Tablica projektu nie ma jeszcze kolumn.');
      }
    });
  }

  function odczytajOs(projekt: string): void {
    void planowanie.harmonogram({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się odczytać osi czasu projektu.', wynik.blad);
        return;
      }
      const harmonogram = wynik.wynik;
      const lista = wykaz('Słupki osi czasu', 'dw-wykaz');
      for (const slupek of harmonogram.bars) {
        const opis = [
          new Date(slupek.startAt).toISOString(),
          new Date(slupek.endAt).toISOString(),
          slupek.critical === true ? 'ścieżka krytyczna' : 'poza ścieżką krytyczną',
        ].join(' → ');
        lista.append(pozycjaWykazu(slupek.title, opis, 'dw').element);
      }
      const poza = harmonogram.unscheduledTaskIds ?? [];
      if (poza.length > 0) {
        lista.append(
          pozycjaWykazu(
            'Poza harmonogramem',
            `${poza.length} zadań bez obu granic czasu: ${poza.join(', ')}`,
            'dw',
          ).element,
        );
      }
      if (harmonogram.bars.length === 0 && poza.length === 0) {
        tresc.pusto('Projekt nie ma zadań z granicami czasu.');
        return;
      }
      tresc.tresc().append(lista);
    });
  }

  function odczytajKalendarz(projekt: string): void {
    void planowanie
      .kalendarz({
        projectId: projekt,
        span: WorkspaceCalendarSpan.Month,
        anchorAt: Date.now(),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się odczytać kalendarza projektu.', wynik.blad);
          return;
        }
        const kalendarz = wynik.wynik;
        if (kalendarz.entries.length === 0) {
          tresc.pusto('Bieżący miesiąc nie ma pozycji kalendarza.');
          return;
        }
        const lista = wykaz('Kalendarz projektu', 'dw-wykaz');
        for (const pozycja of kalendarz.entries) {
          lista.append(
            pozycjaWykazu(pozycja.title, new Date(pozycja.startAt).toISOString(), 'dw').element,
          );
        }
        tresc.tresc().append(lista);
      });
  }

  function odczytaj(): void {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    tresc.ladowanie('Odczyt huba planowania…');
    switch (widok.value) {
      case 'tablica':
        odczytajTablice(projekt);
        return;
      case 'os':
        odczytajOs(projekt);
        return;
      case 'kalendarz':
        odczytajKalendarz(projekt);
        return;
      default:
        odczytajListe(projekt);
    }
  }

  zaloz.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    if (tytul.value.trim() === '') {
      tresc.potwierdzenie('Zadanie bez tytułu nie zostanie założone.', false);
      return;
    }
    void planowanie
      .zalozZadanie({ projectId: projekt, title: tytul.value.trim() })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się założyć zadania.', wynik.blad);
          return;
        }
        wskazaneZadanie.value = wynik.wynik.task.id;
        tytul.value = '';
        tresc.potwierdzenie(`Założono zadanie „${wynik.wynik.task.title}”.`, true);
        odczytaj();
      });
  });

  zmien.addEventListener('click', () => {
    if (wskazaneZadanie.value.trim() === '') {
      tresc.potwierdzenie('Wskaż zadanie, którego stan ma się zmienić.', false);
      return;
    }
    void planowanie
      .zmienZadanie({
        taskId: wskazaneZadanie.value.trim(),
        status: stanZadania.value as WorkspaceTask['status'],
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się zmienić zadania.', wynik.blad);
          return;
        }
        const przesuniete = wynik.wynik.rescheduledTaskIds ?? [];
        tresc.potwierdzenie(
          przesuniete.length === 0
            ? `Zadanie stoi w stanie ${wynik.wynik.task.status}.`
            : `Zadanie stoi w stanie ${wynik.wynik.task.status}; przesunięto ${przesuniete.length} zadań zależnych.`,
          true,
        );
        odczytaj();
      });
  });

  przenies.addEventListener('click', () => {
    if (wskazaneZadanie.value.trim() === '') {
      tresc.potwierdzenie('Wskaż kartę, którą chcesz przenieść.', false);
      return;
    }
    void planowanie
      .przeniesZadanie({
        taskId: wskazaneZadanie.value.trim(),
        status: stanZadania.value as WorkspaceTask['status'],
        boardColumnId: `kolumna-${stanZadania.value}`,
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się przenieść karty.', wynik.blad);
          return;
        }
        tresc.potwierdzenie(
          wynik.wynik.wipExceeded === true
            ? 'Karta przeniesiona; kolumna przekroczyła własną granicę prac w toku.'
            : 'Karta przeniesiona.',
          true,
        );
        odczytaj();
      });
  });

  usun.addEventListener('click', () => {
    if (wskazaneZadanie.value.trim() === '') {
      tresc.potwierdzenie('Wskaż zadanie do usunięcia.', false);
      return;
    }
    void planowanie.usunZadanie({ taskId: wskazaneZadanie.value.trim() }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się usunąć zadania.', wynik.blad);
        return;
      }
      const podrzedne = wynik.wynik.deletedSubtaskIds ?? [];
      tresc.potwierdzenie(
        wynik.wynik.deleted
          ? `Usunięto zadanie wraz z ${podrzedne.length} podzadaniami.`
          : 'Rdzeń nie miał takiego zadania.',
        wynik.wynik.deleted,
      );
      wskazaneZadanie.value = '';
      odczytaj();
    });
  });

  zwiaz.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    if (wskazaneZadanie.value.trim() === '' || drugieZadanie.value.trim() === '') {
      tresc.potwierdzenie('Zależność wymaga wskazania obu zadań.', false);
      return;
    }
    void planowanie
      .zalozZaleznosc({
        projectId: projekt,
        predecessorTaskId: wskazaneZadanie.value.trim(),
        successorTaskId: drugieZadanie.value.trim(),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się związać zadań.', wynik.blad);
          return;
        }
        wskazanaZaleznosc.value = wynik.wynik.dependency.id;
        const przesuniete = wynik.wynik.rescheduledTaskIds ?? [];
        tresc.potwierdzenie(
          `Zadania związane; przesunięto ${przesuniete.length} terminów.`,
          true,
        );
        odczytaj();
      });
  });

  rozwiaz.addEventListener('click', () => {
    if (wskazanaZaleznosc.value.trim() === '') {
      tresc.potwierdzenie('Wskaż zależność, która ma zostać zniesiona.', false);
      return;
    }
    void planowanie
      .zniesZaleznosc({ dependencyId: wskazanaZaleznosc.value.trim() })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się znieść zależności.', wynik.blad);
          return;
        }
        tresc.potwierdzenie(
          wynik.wynik.removed ? 'Zależność zniesiona.' : 'Rdzeń nie miał takiej zależności.',
          wynik.wynik.removed,
        );
        wskazanaZaleznosc.value = '';
        odczytaj();
      });
  });

  wciagnij.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    if (trescIcal.value.trim() === '') {
      tresc.potwierdzenie('Wciągnięcie bez treści pliku nie ma czego wczytać.', false);
      return;
    }
    void planowanie
      .wciagnijKalendarz({ projectId: projekt, content: trescIcal.value })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się wciągnąć kalendarza.', wynik.blad);
          return;
        }
        const powody = wynik.wynik.skippedReasons ?? [];
        tresc.potwierdzenie(
          `Wciągnięto ${wynik.wynik.imported} pozycji; pominięto ${wynik.wynik.skipped}` +
            (powody.length === 0 ? '.' : `: ${powody.join('; ')}`),
          wynik.wynik.imported > 0,
        );
        odczytaj();
      });
  });

  zapiszIcal.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    void planowanie.zapiszKalendarz({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się zapisać kalendarza.', wynik.blad);
        return;
      }
      pobierzPlik(wynik.wynik.fileName, wynik.wynik.content, 'text/calendar');
      tresc.potwierdzenie(`Zapisano ${wynik.wynik.exported} pozycji kalendarza.`, true);
    });
  });

  przestawProjekt.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    void wiedza
      .ustawStanProjektu({
        projectId: projekt,
        status: stanProjektu.value as WorkspaceProjectStatus,
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się ustawić stanu projektu.', wynik.blad);
          return;
        }
        tresc.potwierdzenie(`Projekt stoi w stanie ${wynik.wynik.status}.`, true);
      });
  });

  odswiezOkno.addEventListener('click', () => odczytaj());
  widok.addEventListener('change', () => odczytaj());
  stan.naZmiane(() => odczytaj());

  return {
    element: rama.element,
    odswiez: () => odczytaj(),
  };
}
