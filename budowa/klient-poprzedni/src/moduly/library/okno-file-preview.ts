import { LibraryPreviewKind, type LibraryPreview } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { utworzRameOkna, type RamaOkna } from '../../komponenty/rama-okna';
import { utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { werdyktZPodgladu } from './dostepnosc-tresci';
import { BEZ_KOMENDY_PODGLAD } from './etykiety-biblioteki';
import { utworzPanelAkcji } from './panel-akcji';
import { utworzPasekPodgladu } from './pasek-podgladu';
import type { StanBiblioteki } from './stan-biblioteki';
import { utworzStanOkna } from './stan-okna';
import { utworzSterModulu } from './ster-modulu';
import { utworzSterowanieWidoku } from './sterowanie-widoku';
import { utworzTrescPodgladu } from './tresc-podgladu';
import { przeniesZasoby } from './zapisy-zbiorcze';
import type { ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Okno File Preview aktywuje się po wybraniu pliku w Library Explorer: pokazuje
 * zawartość ze stronicowaniem i przenosi plik do edycji we właściwym module.
 */
export interface OknoPodgladu {
  element: HTMLElement;
  odswiez(): void;
}

/** Obudowa okna: nazwa i rola z wykazu okien operacyjnych oraz objaśnienie dostępne pod znakiem zapytania. */
function utworzRame(): RamaOkna {
  return utworzRameOkna({
    kod: 'file-preview',
    tytul: 'File Preview',
    rola: 'pomocnicze',
    modul: 'Library',
    dodatkiNaglowka: [
      utworzDymekObjasnienia(
        'Podgląd zawartości pliku bez opuszczania modułu. Okno aktywuje się po wybraniu pliku ' +
          'w Library Explorer; „Zamknij" zdejmuje wskazanie, a nie okno komunikacji sesji.',
        { powloka: 'ml-dymek', znak: 'ml-dymek__znak' },
      ),
    ],
  });
}

export function utworzOknoPodgladu(
  stan: StanBiblioteki,
  otoczenie: ZrodloOtoczenia,
): OknoPodgladu {
  const rama = utworzRame();
  const okno = utworzStanOkna();
  const odpowiedz = utworzWierszOdpowiedzi();
  const panel = utworzPanelAkcji(otoczenie, stan, BEZ_KOMENDY_PODGLAD);

  const widok = document.createElement('div');
  widok.className = 'ml-podglad';
  const sterowanie = utworzSterowanieWidoku(widok);

  const modul = utworzSterModulu(
    stan,
    'Ta sama nastawa stoi w czynnościach zbiorczych Library Explorera — wybór jest jeden ' +
      'na cały moduł.',
  );

  let strona = 1;
  let ostatni = '';

  async function pokaz(): Promise<void> {
    const plik = stan.czynny();
    if (plik === null) return;
    okno.ladowanie(`Rdzeń przygotowuje podgląd pliku „${plik.name}" (strona ${strona}).`);
    const wynik = await stan.zrodlo.podglad(plik.id, strona);
    // Ta odpowiedź jest jedynym świadkiem dostępności treści, więc okno ją odkłada do wspólnego stanu.
    stan.zapiszTresc(plik.id, werdyktZPodgladu(wynik.udany ? wynik.wynik?.preview : undefined, wynik.blad));
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Podgląd pliku', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const podglad: LibraryPreview = wynik.wynik.preview;
    strona = podglad.page ?? strona;
    widok.replaceChildren(utworzTrescPodgladu(podglad));
    pasek.ustawStrone(
      podglad.pageCount === undefined ? `Strona ${strona}.` : `Strona ${strona} z ${podglad.pageCount}.`,
    );
    // Odmowa i pustka to dwie różne rzeczy: błąd należy się wyłącznie odpowiedzi, której rdzeń nie dał.
    if (podglad.kind === LibraryPreviewKind.Text && (podglad.text ?? '') === '') {
      okno.puste(
        'Strona bez ani jednego znaku',
        `Rdzeń oddał podgląd tekstowy strony ${strona} i nie ma w nim treści. Przejdź ` +
          'do kolejnej strony albo wskaż inny plik w Library Explorer.',
      );
      return;
    }
    okno.gotowe();
  }

  async function przejdzDoEdycji(): Promise<void> {
    const plik = stan.czynny();
    if (plik === null) {
      odpowiedz.pokaz('Wybierz plik w Library Explorer — podgląd nie ma własnego wejścia.', false);
      return;
    }
    odpowiedz.pokaz(`Przenoszenie pliku „${plik.name}" do modułu docelowego…`, true);
    const wynik = await przeniesZasoby(stan, otoczenie, stan.modulDocelowy(), [plik]);
    odpowiedz.pokaz(wynik.tresc, wynik.powodzenie);
  }

  const pasek = utworzPasekPodgladu(sterowanie.element, {
    wstecz() {
      if (strona <= 1) {
        odpowiedz.pokaz('To jest pierwsza strona podglądu.', false);
        return;
      }
      strona -= 1;
      void pokaz();
    },
    dalej() {
      strona += 1;
      void pokaz();
    },
    otworz: () => void przejdzDoEdycji(),
    zamknij() {
      stan.wskaz(null);
      odpowiedz.pokaz('Podgląd zamknięty — wskazanie pliku zdjęte.', true);
    },
  });

  okno.tresc.append(widok, pasek.element, modul.element);
  rama.cialo.append(okno.element, odpowiedz.element, panel.element);

  return {
    element: rama.element,

    odswiez() {
      // Ster idzie pierwszy, przed każdym wyjściem z funkcji, bo dalsze gałęzie kończą odświeżanie wcześnie.
      modul.odswiez();
      const plik = stan.czynny();
      if (plik === null) {
        widok.replaceChildren();
        okno.puste(
          'Nic nie wybrano',
          'Wskaż plik w Library Explorer — podgląd nie ma własnego wejścia.',
        );
        ostatni = '';
        return;
      }
      if (plik.id === ostatni) return;
      ostatni = plik.id;
      strona = 1;
      sterowanie.zeruj();
      void pokaz();
    },
  };
}
