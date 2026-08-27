import type { LibraryFile } from '../../../../shared/contract';
import { naBase64 } from './biblioteka-pozycja';
import type { StanProjektu } from './stan-projektu';
import type { StanTresci } from './stany-okna';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Czynności okna Project Library sięgające poza sam wykaz — wgranie pliku, udostępnienie
 * zaznaczenia i etykieta zbiorcza — potwierdzają skutek w wykazie zasobów, nie w tym, co wysłało okno.
 */
export interface ZaleznosciCzynnosci {
  zrodlo: ZrodloWorkspace;
  stan: StanProjektu;
  tresc: StanTresci;
  zaznaczone: Set<string>;
  /** Odczytuje wykaz zasobów i oddaje to, co narysował; `null` = odczyt nieudany. */
  odczytaj(): Promise<readonly LibraryFile[] | null>;
}

export interface CzynnosciZbiorcze {
  wgraj(plik: File): void;
  udostepnij(modulDocelowy: string): void;
  oznacz(etykieta: string): void;
}

export function czynnosciZbiorcze(z: ZaleznosciCzynnosci): CzynnosciZbiorcze {
  const { zrodlo, stan, tresc, zaznaczone } = z;

  return {
    wgraj(plik) {
      const projekt = stan.projekt();
      if (projekt === '') {
        tresc.blad('Wgranie bez projektu nie ma zakresu — wskaż projekt na pulpicie.');
        return;
      }
      tresc.ladowanie(`Wgrywanie pliku ${plik.name}…`);
      void plik
        .arrayBuffer()
        .then((zawartosc) => zrodlo.wgrajPlik(projekt, plik.name, naBase64(zawartosc)))
        .then(async (wynik) => {
          if (!wynik.udany || wynik.wynik === undefined) {
            tresc.blad(`Rdzeń nie przyjął pliku ${plik.name}.`, wynik.blad);
            return;
          }
          const oddany = wynik.wynik;
          const wykaz = await z.odczytaj();
          if (wykaz !== null && !wykaz.some((pozycja) => pozycja.id === oddany.id)) {
            tresc.potwierdzenie(
              `Rdzeń zapisał plik ${oddany.name} (id ${oddany.id}, projekt ` +
                `${oddany.projectId ?? 'nie podany'}), ale nie ma go w wykazie zasobów ` +
                `projektu ${projekt} — wgranie DO PROJEKTU się nie odbyło.`,
              false,
            );
            return;
          }
          tresc.potwierdzenie(
            `Rdzeń zapisał plik ${oddany.name} w projekcie ${oddany.projectId ?? 'nie podanym'}.`,
            true,
          );
        });
    },

    udostepnij(modulDocelowy) {
      const okno = stan.oknoRozmowy();
      if (zaznaczone.size === 0) {
        tresc.potwierdzenie('Zaznacz przynajmniej jeden plik, aby go udostępnić.', false);
        return;
      }
      if (okno === '') {
        // Powód bierze się ze stanu, a ten z odpowiedzi rdzenia na wykaz okien, rozróżniając odmowę od braku.
        tresc.blad(`Przeniesienie kontekstu wymaga okna rozmowy modułu. ${stan.powodBrakuOkna()}`);
        return;
      }
      tresc.ladowanie('Przenoszenie materiału do modułu docelowego…');
      void zrodlo
        .udostepnij(okno, modulDocelowy, stan.projekt(), [...zaznaczone])
        .then((wynik) => {
          if (!wynik.udany || wynik.wynik === undefined) {
            tresc.blad('Rdzeń nie przeniósł materiału do modułu docelowego.', wynik.blad);
            return;
          }
          void z.odczytaj();
          // Odpowiedź przeniesienia niesie znacznik przeniesienia i okno docelowe, oba stamtąd bierze zdanie.
          const odpowiedz = wynik.wynik;
          if (odpowiedz.transferred !== true) {
            tresc.potwierdzenie(
              'Rdzeń przyjął wywołanie, ale nie potwierdził przeniesienia ' +
                '(context.transfer oddał transferred = fałsz) — materiał został na miejscu.',
              false,
            );
            return;
          }
          tresc.potwierdzenie(
            `Rdzeń przeniósł materiał do okna ${odpowiedz.window.id} ` +
              `modułu ${odpowiedz.window.moduleId}.`,
            true,
          );
        });
    },

    oznacz(etykieta) {
      const etykiety = etykieta
        .split(',')
        .map((pozycja) => pozycja.trim())
        .filter((pozycja) => pozycja !== '');
      if (zaznaczone.size === 0 || etykiety.length === 0) {
        tresc.potwierdzenie('Etykieta zbiorcza wymaga zaznaczenia i treści etykiety.', false);
        return;
      }
      tresc.ladowanie(`Nadawanie etykiety ${zaznaczone.size} plikom…`);
      // Kontrakt nadaje etykietę po pliku, więc czynność zbiorcza jest pętlą, a odmowa wychodzi na wierzch.
      void Promise.all(
        [...zaznaczone].map((identyfikator) => zrodlo.nadajEtykiete(identyfikator, etykiety)),
      ).then((wyniki) => {
        const nieudane = wyniki.filter((wynik) => !wynik.udany);
        if (nieudane.length > 0) {
          tresc.blad(
            `Etykieta nie została nadana ${nieudane.length} z ${wyniki.length} plików.`,
            nieudane[0]?.blad,
          );
          return;
        }
        void z.odczytaj();
        // Rdzeń oddaje plik po zapisie wraz z etykietami — sprawdzenie idzie po tej odpowiedzi, nie po polu.
        const bezEtykiety = wyniki.filter(
          (wynik) =>
            wynik.wynik === undefined ||
            !etykiety.every((pozycja) => (wynik.wynik?.tags ?? []).includes(pozycja)),
        );
        if (bezEtykiety.length > 0) {
          tresc.potwierdzenie(
            `Rdzeń przyjął wywołanie, ale w ${bezEtykiety.length} z ${wyniki.length} ` +
              `plików nie oddał etykiet ${etykiety.join(', ')} — nadanie się nie odbyło.`,
            false,
          );
          return;
        }
        tresc.potwierdzenie(
          `Rdzeń zapisał etykiety ${etykiety.join(', ')} na ${wyniki.length} plikach.`,
          true,
        );
      });
    },
  };
}
