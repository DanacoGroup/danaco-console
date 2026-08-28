import { WorkspaceEntityKind } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleTresci,
  pozycjaWykazu,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { CzynnosciWiedzy } from './czynnosci-wiedzy';
import type { StanProjektu } from './stan-projektu';
import { utworzStanTresci } from './stany-okna';

/**
 * Współpraca i materiały projektu gromadzą czynności przekrojowe: wyszukiwanie po słowach, oś
 * czasu, komentarze, wydobycie treści, duplikaty i historię instrukcji.
 */
export interface OknoWspolpracy {
  element: HTMLElement;
  odswiez(): void;
}

/** Rodzaje bytu komentowanego to wykaz wyliczony z kontraktu, nie przepisany, dający listę wyboru rodzaju bytu w formularzu. */
const RODZAJE_BYTU: ReadonlyArray<readonly [string, string]> = Object.values(
  WorkspaceEntityKind,
).map((rodzaj) => [rodzaj, `byt: ${rodzaj}`]);

export function utworzOknoWspolpracy(
  wiedza: CzynnosciWiedzy,
  stan: StanProjektu,
): OknoWspolpracy {
  const rama = utworzRameOkna({
    tytul: 'Współpraca i materiały',
    kod: 'project-collaboration',
    rola: 'zarządca',
    przeznaczenie:
      'Wyszukiwanie w projekcie, oś czasu aktywności, komentarze, wydobycie treści plików i historia instrukcji.',
    przedrostek: 'dw',
  });
  const tresc = utworzStanTresci();

  const fraza = pole('Fraza wyszukiwania', 'słowo szukane w projekcie');
  const plik = pole('Plik projektu', 'ścieżka pliku w bibliotece projektu');
  const drugiPlik = pole('Pliki scalane', 'ścieżki rozdzielone przecinkiem');
  const byt = pole('Byt komentowany', 'identyfikator zadania, notatki albo pliku');
  const rodzajBytu = wybor('Rodzaj bytu', RODZAJE_BYTU);
  const trescKomentarza = poleTresci('Treść komentarza', 3);
  const wskazanyKomentarz = pole('Komentarz wskazany', 'identyfikator komentarza');
  const wersja = pole('Wersja instrukcji', 'identyfikator wersji z wykazu');
  const ekspert = pole('Ekspert', 'identyfikator eksperta do odłączenia');

  const szukaj = przycisk('Szukaj w projekcie', 'dn-btn dn-btn--atrament');
  const osCzasu = przycisk('Oś czasu aktywności');
  const dolozKomentarz = przycisk('Dołóż komentarz');
  const wykazKomentarzy = przycisk('Wykaz komentarzy');
  const usunKomentarz = przycisk('Usuń komentarz');
  const wydobadz = przycisk('Wydobądź treść pliku');
  const duplikaty = przycisk('Znajdź duplikaty');
  const scal = przycisk('Scal duplikaty');
  const wersje = przycisk('Wersje instrukcji');
  const przywroc = przycisk('Przywróć wersję');
  const odlacz = przycisk('Odłącz eksperta');

  rama.akcje.append(
    szukaj,
    osCzasu,
    dolozKomentarz,
    wykazKomentarzy,
    usunKomentarz,
    wydobadz,
    duplikaty,
    scal,
    wersje,
    przywroc,
    odlacz,
  );
  rama.narzedzia.append(rodzajBytu);
  rama.cialo.append(
    wiersz('Fraza', fraza, {
      klasa: 'dw-wiersz',
      objasnienie: 'Wyszukiwanie po słowach obejmuje zadania, notatki, pliki, pamięć i instrukcje.',
    }),
    wiersz('Plik projektu', plik, {
      klasa: 'dw-wiersz',
      objasnienie:
        'Wydobycie czyta warstwę tekstową dokumentu, a przy jej braku rozpoznaje pismo ze skanu.',
    }),
    wiersz('Pliki scalane', drugiPlik, {
      klasa: 'dw-wiersz',
      objasnienie: 'Scalenie usuwa wskazane pliki; plik o innej treści zostaje nietknięty.',
    }),
    wiersz('Byt komentowany', byt, { klasa: 'dw-wiersz' }),
    wiersz('Treść komentarza', trescKomentarza, {
      klasa: 'dw-wiersz',
      objasnienie: 'Znak małpy w treści przywołuje byt projektu; przywołania wracają w odpowiedzi.',
    }),
    wiersz('Komentarz wskazany', wskazanyKomentarz, { klasa: 'dw-wiersz' }),
    wiersz('Wersja instrukcji', wersja, {
      klasa: 'dw-wiersz',
      objasnienie: 'Przywrócenie zakłada wersję nową o treści wskazanej; historia zostaje.',
    }),
    wiersz('Ekspert', ekspert, {
      klasa: 'dw-wiersz',
      objasnienie: 'Odłączenie znosi przypisanie do projektu, nie usuwa eksperta z modułu Agents.',
    }),
    tresc.element,
  );

  function projektAlboOstrzez(): string {
    const projekt = stan.projekt();
    if (projekt === '') {
      tresc.blad('Czynności przekrojowe bez projektu nie mają zakresu — wskaż projekt.');
      return '';
    }
    return projekt;
  }

  szukaj.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    if (fraza.value.trim() === '') {
      tresc.potwierdzenie('Wyszukiwanie bez frazy nie ma czego szukać.', false);
      return;
    }
    tresc.ladowanie('Wyszukiwanie w projekcie…');
    void wiedza.szukaj({ projectId: projekt, query: fraza.value.trim() }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się przeszukać projektu.', wynik.blad);
        return;
      }
      if (wynik.wynik.length === 0) {
        tresc.pusto('Żaden byt projektu nie niesie tej frazy.');
        return;
      }
      const lista = wykaz('Trafienia wyszukiwania', 'dw-wykaz');
      for (const trafienie of wynik.wynik) {
        lista.append(
          pozycjaWykazu(
            `${trafienie.title} (${trafienie.kind})`,
            trafienie.snippet ?? 'trafienie w nazwę bytu',
            'dw',
          ).element,
        );
      }
      tresc.tresc().append(lista);
    });
  });

  osCzasu.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    tresc.ladowanie('Odczyt osi czasu projektu…');
    void wiedza.osCzasu({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się odczytać osi czasu.', wynik.blad);
        return;
      }
      if (wynik.wynik.length === 0) {
        tresc.pusto('Oś czasu projektu jest pusta — w projekcie nic się jeszcze nie wydarzyło.');
        return;
      }
      const lista = wykaz('Oś czasu aktywności', 'dw-wykaz');
      for (const zdarzenie of wynik.wynik) {
        lista.append(
          pozycjaWykazu(
            zdarzenie.summary,
            `${zdarzenie.kind} · ${new Date(zdarzenie.occurredAt).toISOString()}`,
            'dw',
          ).element,
        );
      }
      tresc.tresc().append(lista);
    });
  });

  dolozKomentarz.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    if (byt.value.trim() === '' || trescKomentarza.value.trim() === '') {
      tresc.potwierdzenie('Komentarz wymaga bytu i treści.', false);
      return;
    }
    void wiedza
      .dolozKomentarz({
        projectId: projekt,
        targetKind: rodzajBytu.value as WorkspaceEntityKind,
        targetId: byt.value.trim(),
        content: trescKomentarza.value,
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się dołożyć komentarza.', wynik.blad);
          return;
        }
        wskazanyKomentarz.value = wynik.wynik.comment.id;
        const przywolania = wynik.wynik.mentionedIds ?? [];
        tresc.potwierdzenie(
          przywolania.length === 0
            ? 'Komentarz zapisany.'
            : `Komentarz zapisany; przywołano: ${przywolania.join(', ')}.`,
          true,
        );
        trescKomentarza.value = '';
      });
  });

  wykazKomentarzy.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    tresc.ladowanie('Odczyt komentarzy projektu…');
    const zadanie =
      byt.value.trim() === ''
        ? { projectId: projekt }
        : { projectId: projekt, targetId: byt.value.trim() };
    void wiedza.komentarze(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się odczytać komentarzy.', wynik.blad);
        return;
      }
      if (wynik.wynik.length === 0) {
        tresc.pusto('Wskazany zakres nie ma komentarzy.');
        return;
      }
      const lista = wykaz('Komentarze projektu', 'dw-wykaz');
      for (const komentarz of wynik.wynik) {
        lista.append(
          pozycjaWykazu(
            komentarz.content,
            `${komentarz.targetKind} ${komentarz.targetId} · ${komentarz.authorKind}`,
            'dw',
          ).element,
        );
      }
      tresc.tresc().append(lista);
    });
  });

  usunKomentarz.addEventListener('click', () => {
    if (wskazanyKomentarz.value.trim() === '') {
      tresc.potwierdzenie('Wskaż komentarz do usunięcia.', false);
      return;
    }
    void wiedza
      .usunKomentarz({ commentId: wskazanyKomentarz.value.trim() })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się usunąć komentarza.', wynik.blad);
          return;
        }
        const usuniete = wynik.wynik.deletedCommentIds ?? [];
        tresc.potwierdzenie(
          wynik.wynik.deleted
            ? `Usunięto ${usuniete.length} komentarzy wątku.`
            : 'Rdzeń nie miał takiego komentarza.',
          wynik.wynik.deleted,
        );
        wskazanyKomentarz.value = '';
      });
  });

  wydobadz.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    if (plik.value.trim() === '') {
      tresc.potwierdzenie('Wskaż plik, z którego ma zostać wydobyta treść.', false);
      return;
    }
    tresc.ladowanie('Wydobycie treści pliku…');
    void wiedza
      .wydobadzTekst({ projectId: projekt, fileId: plik.value.trim() })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się wydobyć treści pliku.', wynik.blad);
          return;
        }
        const wyciag = wynik.wynik.extraction;
        const lista = wykaz('Wynik wydobycia', 'dw-wykaz');
        lista.append(
          pozycjaWykazu('Sposób', wyciag.method, 'dw').element,
          pozycjaWykazu('Znaków', String(wyciag.characterCount), 'dw').element,
          pozycjaWykazu('Początek treści', wyciag.excerpt ?? 'brak podglądu', 'dw').element,
        );
        tresc.tresc().append(lista);
        tresc.potwierdzenie('Treść pliku weszła do wskaźnika wyszukiwania projektu.', true);
      });
  });

  duplikaty.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    tresc.ladowanie('Szukanie plików o treści identycznej…');
    void wiedza.duplikaty({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się wykryć duplikatów.', wynik.blad);
        return;
      }
      if (wynik.wynik.groups.length === 0) {
        tresc.pusto('Biblioteka projektu nie ma plików o treści identycznej.');
        return;
      }
      const lista = wykaz('Grupy duplikatów', 'dw-wykaz');
      for (const grupa of wynik.wynik.groups) {
        const nazwy = grupa.files.map((wpis) => wpis.id).join(', ');
        const { element, akcje } = pozycjaWykazu(
          `${grupa.files.length} plików o tej samej treści`,
          nazwy,
          'dw',
        );
        const wskaz = przycisk('Wskaż do scalenia', 'dn-btn dn-btn--zarys');
        wskaz.addEventListener('click', () => {
          plik.value = grupa.suggestedKeepFileId ?? grupa.files[0]?.id ?? '';
          drugiPlik.value = grupa.files
            .map((wpis) => wpis.id)
            .filter((kod) => kod !== plik.value)
            .join(', ');
          tresc.potwierdzenie('Wskazano grupę do scalenia.', true);
        });
        akcje.append(wskaz);
        lista.append(element);
      }
      tresc.tresc().append(lista);
    });
  });

  scal.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    const scalane = drugiPlik.value
      .split(',')
      .map((kod) => kod.trim())
      .filter((kod) => kod !== '');
    if (plik.value.trim() === '' || scalane.length === 0) {
      tresc.potwierdzenie('Scalenie wymaga pliku zachowywanego i plików scalanych.', false);
      return;
    }
    void wiedza
      .scalDuplikaty({
        projectId: projekt,
        keepFileId: plik.value.trim(),
        mergedFileIds: scalane,
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się scalić duplikatów.', wynik.blad);
          return;
        }
        tresc.potwierdzenie(
          `Scalono ${wynik.wynik.merged} plików; odzyskano ${wynik.wynik.reclaimedBytes ?? 0} bajtów.`,
          wynik.wynik.merged > 0,
        );
        drugiPlik.value = '';
      });
  });

  wersje.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    tresc.ladowanie('Odczyt historii instrukcji…');
    void wiedza.wersjeInstrukcji({ projectId: projekt }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Nie udało się odczytać wersji instrukcji.', wynik.blad);
        return;
      }
      if (wynik.wynik.length === 0) {
        tresc.pusto('Instrukcje projektu nie mają jeszcze historii — pierwszy zapis ją założy.');
        return;
      }
      wersja.value = wynik.wynik[0]?.id ?? '';
      const lista = wykaz('Wersje instrukcji projektu', 'dw-wykaz');
      for (const pozycja of wynik.wynik) {
        const { element, akcje } = pozycjaWykazu(
          new Date(pozycja.createdAt).toISOString(),
          pozycja.content.slice(0, 120),
          'dw',
        );
        const wskaz = przycisk('Wskaż', 'dn-btn dn-btn--zarys');
        wskaz.addEventListener('click', () => {
          wersja.value = pozycja.id;
          tresc.potwierdzenie('Wskazano wersję instrukcji.', true);
        });
        akcje.append(wskaz);
        lista.append(element);
      }
      tresc.tresc().append(lista);
    });
  });

  przywroc.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    if (wersja.value.trim() === '') {
      tresc.potwierdzenie('Wskaż wersję instrukcji do przywrócenia.', false);
      return;
    }
    void wiedza
      .przywrocInstrukcje({ projectId: projekt, versionId: wersja.value.trim() })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się przywrócić wersji instrukcji.', wynik.blad);
          return;
        }
        tresc.potwierdzenie(
          `Przywrócono treść wersji; historia liczy teraz wersję ${wynik.wynik.version.id}.`,
          true,
        );
      });
  });

  odlacz.addEventListener('click', () => {
    const projekt = projektAlboOstrzez();
    if (projekt === '') return;
    if (ekspert.value.trim() === '') {
      tresc.potwierdzenie('Wskaż eksperta do odłączenia.', false);
      return;
    }
    void wiedza
      .odlaczAgenta({ projectId: projekt, agentId: ekspert.value.trim() })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się odłączyć eksperta.', wynik.blad);
          return;
        }
        tresc.potwierdzenie(
          wynik.wynik.unassigned
            ? 'Ekspert odłączony od projektu.'
            : 'Ekspert nie był przypisany do tego projektu.',
          wynik.wynik.unassigned,
        );
      });
  });

  function odswiez(): void {
    const projekt = stan.projekt();
    if (projekt === '') {
      tresc.pusto('Wskaż projekt, aby korzystać z czynności przekrojowych.');
      return;
    }
    tresc.pusto('Wybierz czynność: wyszukiwanie, oś czasu, komentarze, materiały albo instrukcje.');
  }

  stan.naZmiane(() => odswiez());

  return { element: rama.element, odswiez };
}
