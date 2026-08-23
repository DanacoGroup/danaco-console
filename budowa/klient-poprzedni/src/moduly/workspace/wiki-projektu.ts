import type { WorkspaceNote } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleTresci,
  pozycjaWykazu,
  przyciskAkcji as przycisk,
  wiersz,
  wykaz,
} from '../../modele/kontrolki-formularza';
import {
  glebokoscWezla,
  nazwyOdnosnikow,
  type CzynnosciWiedzy,
} from './czynnosci-wiedzy';
import type { StanProjektu } from './stan-projektu';
import { utworzStanTresci } from './stany-okna';

/**
 * Notatki i wiki projektu — zakładka „Notatki i wiki" oraz „Tablica" okna
 * Project Library, wydzielona w osobne okno modułu.
 *
 * Strona jest notatką: hierarchię daje wskazanie strony nadrzędnej, a odnośniki
 * zapisem `[[nazwa]]` w treści. Odnośnik do strony jeszcze niezałożonej jest
 * stanem poprawnym wiki — okno pokazuje takie nazwy wprost, zamiast milczeć
 * i zostawiać Operatora z martwym odnośnikiem do odkrycia po kliknięciu.
 *
 * Tablica wizualna stoi w tym samym oknie, bo niesie ten sam materiał w innej
 * formie: scena płótna jest zapisem JSON i okno oddaje ją Operatorowi wprost,
 * zamiast udawać, że rysuje płótno, którego nie rysuje.
 */
export interface OknoWikiProjektu {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoWikiProjektu(
  wiedza: CzynnosciWiedzy,
  stan: StanProjektu,
): OknoWikiProjektu {
  const rama = utworzRameOkna({
    tytul: 'Notatki i wiki',
    kod: 'project-wiki',
    rola: 'zarządca',
    przeznaczenie:
      'Strony wiki projektu wraz z drzewem, odnośnikami wstecznymi, grafem powiązań i tablicą wizualną.',
    przedrostek: 'dw',
  });
  const tresc = utworzStanTresci();

  const tytul = pole('Tytuł notatki', 'np. Protokół spotkania');
  const trescNotatki = poleTresci('Treść notatki (Markdown)', 5);
  const nadrzedna = pole('Strona nadrzędna', 'identyfikator strony nadrzędnej');
  const wskazana = pole('Notatka wskazana', 'identyfikator notatki z wykazu');
  const scena = poleTresci('Scena tablicy wizualnej (JSON)', 3);

  const zapisz = przycisk('Zapisz notatkę', 'dn-btn dn-btn--atrament');
  const otworz = przycisk('Otwórz notatkę');
  const wykazNotatek = przycisk('Wykaz notatek');
  const drzewo = przycisk('Drzewo stron');
  const wsteczne = przycisk('Co linkuje tutaj');
  const graf = przycisk('Graf wiedzy');
  const usun = przycisk('Usuń notatkę');
  const wczytajTablice = przycisk('Wczytaj tablicę');
  const zapiszTablice = przycisk('Zapisz tablicę');

  rama.akcje.append(
    zapisz,
    otworz,
    wykazNotatek,
    drzewo,
    wsteczne,
    graf,
    usun,
    wczytajTablice,
    zapiszTablice,
  );
  rama.cialo.append(
    wiersz('Tytuł', tytul, {
      klasa: 'dw-wiersz',
      objasnienie: 'Tytuł jest nazwą, którą wskazują odnośniki [[nazwa]] innych stron.',
    }),
    wiersz('Treść', trescNotatki, {
      klasa: 'dw-wiersz',
      objasnienie: 'Zapis [[nazwa]] wiąże stronę ze stroną; nagłówki dają spis treści.',
    }),
    wiersz('Strona nadrzędna', nadrzedna, {
      klasa: 'dw-wiersz',
      objasnienie: 'Puste znaczy stronę korzenia drzewa.',
    }),
    wiersz('Notatka wskazana', wskazana, {
      klasa: 'dw-wiersz',
      objasnienie: 'Kliknięcie „Wskaż” przy pozycji wykazu wpisuje tu jej identyfikator.',
    }),
    wiersz('Scena tablicy', scena, {
      klasa: 'dw-wiersz',
      objasnienie: 'Zapis JSON płótna: kartki, strzałki i grupy. Rdzeń sceny nie rozbiera.',
    }),
    tresc.element,
  );

  function projektAlboOstrzez(): string {
    const projekt = stan.projekt();
    if (projekt === '') {
      tresc.blad('Wiki bez projektu nie ma zakresu — wskaż projekt na pulpicie.');
      return '';
    }
    return projekt;
  }

  function pozycjaNotatki(notatka: WorkspaceNote): HTMLElement {
    const { element, akcje } = pozycjaWykazu(
      notatka.title,
      notatka.path ?? 'strona korzenia',
      'dw',
    );
    const wskaz = przycisk('Wskaż', 'dn-btn dn-btn--zarys');
    wskaz.addEventListener('click', () => {
      wskazana.value = notatka.id;
      tresc.potwierdzenie(`Wskazano notatkę „${notatka.title}”.`, true);
    });
    akcje.append(wskaz);
    return element;
  }

  function odczytaj(): void {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    tresc.ladowanie('Odczyt notatek projektu…');
    void wiedza.notatki({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się odczytać notatek projektu.', wynik.blad);
        return;
      }
      rama.ustawZnacznik(`${wynik.wynik.length} stron`);
      if (wynik.wynik.length === 0) {
        tresc.pusto('Projekt nie ma jeszcze notatek. Pierwsza strona powstaje zapisem powyżej.');
        return;
      }
      const lista = wykaz('Notatki projektu', 'dw-wykaz');
      for (const notatka of wynik.wynik) lista.append(pozycjaNotatki(notatka));
      tresc.tresc().append(lista);
    });
  }

  zapisz.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    if (tytul.value.trim() === '') {
      tresc.potwierdzenie('Notatka bez tytułu nie zostanie zapisana.', false);
      return;
    }
    const zapowiedziane = nazwyOdnosnikow(trescNotatki.value);
    void wiedza
      .zapiszNotatke({
        projectId: projekt,
        noteId: wskazana.value.trim(),
        title: tytul.value.trim(),
        content: trescNotatki.value,
        parentNoteId: nadrzedna.value.trim(),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się zapisać notatki.', wynik.blad);
          return;
        }
        wskazana.value = wynik.wynik.note.id;
        const brakujace = wynik.wynik.missingNames ?? [];
        tresc.potwierdzenie(
          brakujace.length === 0
            ? `Zapisano notatkę z ${zapowiedziane.length} odnośnikami.`
            : `Zapisano notatkę; bez strony zostają nazwy: ${brakujace.join(', ')}.`,
          true,
        );
        odczytaj();
      });
  });

  otworz.addEventListener('click', () => {
    if (wskazana.value.trim() === '') {
      tresc.potwierdzenie('Wskaż notatkę, którą chcesz otworzyć.', false);
      return;
    }
    void wiedza.notatka({ noteId: wskazana.value.trim() }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się otworzyć notatki.', wynik.blad);
        return;
      }
      tytul.value = wynik.wynik.title;
      trescNotatki.value = wynik.wynik.content;
      nadrzedna.value = wynik.wynik.parentNoteId ?? '';
      const naglowki = wynik.wynik.headings ?? [];
      tresc.potwierdzenie(
        `Otwarto notatkę „${wynik.wynik.title}” (${naglowki.length} nagłówków).`,
        true,
      );
    });
  });

  wykazNotatek.addEventListener('click', () => odczytaj());

  drzewo.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    tresc.ladowanie('Odczyt drzewa stron…');
    void wiedza.drzewoNotatek({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się odczytać drzewa stron.', wynik.blad);
        return;
      }
      if (wynik.wynik.length === 0) {
        tresc.pusto('Drzewo stron jest puste.');
        return;
      }
      const lista = wykaz('Drzewo stron projektu', 'dw-wykaz');
      for (const wezel of wynik.wynik) {
        const wciecie = '· '.repeat(glebokoscWezla(wynik.wynik, wezel));
        lista.append(
          pozycjaWykazu(
            `${wciecie}${wezel.title}`,
            `${wezel.childCount} stron podrzędnych`,
            'dw',
          ).element,
        );
      }
      tresc.tresc().append(lista);
    });
  });

  wsteczne.addEventListener('click', () => {
    if (wskazana.value.trim() === '') {
      tresc.potwierdzenie('Wskaż stronę, dla której szukamy odnośników wstecznych.', false);
      return;
    }
    tresc.ladowanie('Odczyt odnośników wstecznych…');
    void wiedza.odnosnikiWsteczne({ noteId: wskazana.value.trim() }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się odczytać odnośników wstecznych.', wynik.blad);
        return;
      }
      if (wynik.wynik.length === 0) {
        tresc.pusto('Do tej strony nie linkuje żadna inna strona projektu.');
        return;
      }
      const lista = wykaz('Co linkuje tutaj', 'dw-wykaz');
      for (const odnosnik of wynik.wynik) {
        lista.append(
          pozycjaWykazu(odnosnik.sourceTitle, odnosnik.context ?? odnosnik.targetName, 'dw')
            .element,
        );
      }
      tresc.tresc().append(lista);
    });
  });

  graf.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    tresc.ladowanie('Odczyt grafu wiedzy…');
    void wiedza.grafWiedzy({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się odczytać grafu wiedzy.', wynik.blad);
        return;
      }
      const graf = wynik.wynik;
      if (graf.nodes.length === 0) {
        tresc.pusto('Graf wiedzy projektu jest pusty.');
        return;
      }
      const lista = wykaz('Graf wiedzy projektu', 'dw-wykaz');
      lista.append(
        pozycjaWykazu(
          'Rozmiar grafu',
          `${graf.nodes.length} węzłów, ${graf.edges.length} krawędzi` +
            (graf.truncated === true ? ' (obraz przycięty granicą wielkości)' : ''),
          'dw',
        ).element,
      );
      for (const krawedz of graf.edges) {
        lista.append(
          pozycjaWykazu(
            `${krawedz.sourceKind} → ${krawedz.targetKind}`,
            `${krawedz.sourceId} → ${krawedz.targetId} (${krawedz.kind})`,
            'dw',
          ).element,
        );
      }
      tresc.tresc().append(lista);
    });
  });

  usun.addEventListener('click', () => {
    if (wskazana.value.trim() === '') {
      tresc.potwierdzenie('Wskaż notatkę do usunięcia.', false);
      return;
    }
    void wiedza.usunNotatke({ noteId: wskazana.value.trim() }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się usunąć notatki.', wynik.blad);
        return;
      }
      const osierocone = wynik.wynik.orphanedBacklinkCount ?? 0;
      tresc.potwierdzenie(
        wynik.wynik.deleted
          ? `Usunięto notatkę; bez strony zostaje ${osierocone} odnośników.`
          : 'Rdzeń nie miał takiej notatki.',
        wynik.wynik.deleted,
      );
      wskazana.value = '';
      odczytaj();
    });
  });

  wczytajTablice.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    void wiedza.kanwa({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się wczytać tablicy wizualnej.', wynik.blad);
        return;
      }
      const kanwa = wynik.wynik.canvas;
      if (kanwa === undefined) {
        tresc.potwierdzenie('Projekt nie ma jeszcze tablicy wizualnej — zapis założy pierwszą.', true);
        return;
      }
      scena.value = JSON.stringify(kanwa.scene, null, 2);
      tresc.potwierdzenie(`Wczytano tablicę „${kanwa.name}”.`, true);
    });
  });

  zapiszTablice.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    let sceneJson: unknown;
    try {
      sceneJson = JSON.parse(scena.value === '' ? '{}' : scena.value);
    } catch (powod) {
      tresc.potwierdzenie(`Scena nie jest poprawnym zapisem JSON: ${String(powod)}`, false);
      return;
    }
    void wiedza
      .zapiszKanwe({ projectId: projekt, scene: sceneJson })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się zapisać tablicy wizualnej.', wynik.blad);
          return;
        }
        tresc.potwierdzenie(`Zapisano tablicę „${wynik.wynik.canvas.name}”.`, true);
      });
  });

  stan.naZmiane(() => odczytaj());

  return { element: rama.element, odswiez: odczytaj };
}
