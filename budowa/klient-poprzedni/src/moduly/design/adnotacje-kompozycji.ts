import { Command, type DesignAnnotation, type DesignPresence } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleLogiczne,
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { StanKompozycji } from './stan-kompozycji';
import type { StanDesignu } from './stan-designu';

/**
 * Adnotacje i obecność na kompozycji Design Board. Adnotacja ma własny wiersz
 * i własny czas, przeżywa zapis układu i nie ginie przy przesunięciu warstwy.
 * Obecność jest ulotna: zgłoszenie nie zapisuje się w bazie.
 */
export interface AdnotacjeKompozycji {
  element: HTMLElement;
  /** Zleca odczyt adnotacji kompozycji wskazanej na kanwie. */
  wczytaj(): Promise<void>;
  /** Odpina subskrypcję zdarzenia obecności. */
  rozlacz(): void;
}

export function utworzAdnotacjeKompozycji(
  stan: StanDesignu,
  kompozycja: StanKompozycji,
): AdnotacjeKompozycji {
  let zbior: readonly DesignAnnotation[] = [];

  const tresc = poleTekstowe({
    etykieta: 'Treść adnotacji',
    podpowiedz: 'np. margines nagłówka za wąski',
    opis: 'Pole text żądania design.annotation.set.',
  });
  const wybor = poleWyboru(
    {
      etykieta: 'Adnotacja',
      opis:
        'Wskazana adnotacja jest celem zmiany (pole annotationId) albo adnotacją nadrzędną ' +
        'odpowiedzi (pole parentId). Pozycja „nowy wątek" zakłada uwagę bez nadrzędnej.',
    },
    [{ wartosc: '', etykieta: 'nowy wątek' }],
  );
  const doWarstwy = poleLogiczne({
    etykieta: 'Przypnij do warstwy zaznaczonej',
    opis:
      'Pole layerId. Bez zaznaczenia adnotacja przypina się do kompozycji — ' +
      'tak stanowi kontrakt dla pola pominiętego.',
  });
  const tylkoOtwarte = poleLogiczne({
    etykieta: 'Tylko wątki niezamknięte',
    opis: 'Pole openOnly żądania design.annotation.list.',
  });

  const zaloz = przycisk('Zapisz adnotację', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedzWWatku = przycisk('Odpowiedz w wątku', 'dn-btn dn-btn--sm dn-btn--zarys');
  const zamknij = przycisk('Zamknij wątek', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odczytaj = przycisk('Odczytaj adnotacje', 'dn-btn dn-btn--sm dn-btn--zarys');
  const zglos = przycisk('Zgłoś obecność na kompozycji', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odejdz = przycisk('Zdejmij obecność', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const pasek = document.createElement('div');
  pasek.className = 'md-czynnosci__pasek';
  pasek.append(zaloz, odpowiedzWWatku, zamknij, odczytaj, zglos, odejdz);

  const wykaz = document.createElement('ul');
  wykaz.className = 'md-adnotacje__wykaz';

  const kursory = document.createElement('p');
  kursory.className = 'md-adnotacje__kursory';
  kursory.textContent = 'Obecność na kompozycji nie była jeszcze zgłaszana.';

  const element = document.createElement('div');
  element.className = 'md-adnotacje';
  element.append(
    tresc.element,
    wybor.element,
    doWarstwy.element,
    tylkoOtwarte.element,
    pasek,
    wykaz,
    kursory,
    odpowiedz.element,
  );

  zaloz.addEventListener('click', () => void zapisz(false, null));
  odpowiedzWWatku.addEventListener('click', () => void zapisz(true, null));
  zamknij.addEventListener('click', () => void zamknijWatek());
  odczytaj.addEventListener('click', () => void wczytaj());
  zglos.addEventListener('click', () => void zglosObecnosc(false));
  odejdz.addEventListener('click', () => void zglosObecnosc(true));

  // Zdarzenie niesie komplet obecnych, więc panel je podmienia, zamiast doliczać stan przyrostów.
  const odsubskrybuj = stan.zrodlo.naObecnosc((zdarzenie) => {
    if (zdarzenie.boardId !== kompozycja.idKompozycji()) return;
    pokazObecnych(zdarzenie.participants);
  });

  function pokazObecnych(obecni: readonly DesignPresence[]): void {
    kursory.textContent =
      obecni.length === 0
        ? 'Na tej kompozycji nie ma w tej chwili nikogo.'
        : `Obecni na kompozycji: ${obecni
            .map((osoba) => `${osoba.label ?? osoba.clientId}`)
            .join(', ')}.`;
  }

  function bezKompozycji(komenda: string): boolean {
    if (kompozycja.idKompozycji() !== '') return false;
    odpowiedz.pokaz(
      `Komenda ${komenda} wskazuje kompozycję identyfikatorem nadanym przez rdzeń — ` +
        'najpierw zapisz kompozycję przyciskiem „Zapisz kompozycję w rdzeniu".',
      false,
    );
    return true;
  }

  function pokazAdnotacje(): void {
    ustawPozycje(wybor.kontrolka, [
      { wartosc: '', etykieta: 'nowy wątek' },
      ...zbior.map((adnotacja) => ({
        wartosc: adnotacja.id,
        etykieta: `${adnotacja.resolved === true ? '✓ ' : ''}${adnotacja.text}`,
      })),
    ]);
    wykaz.replaceChildren(
      ...zbior.map((adnotacja) => {
        const wiersz = document.createElement('li');
        wiersz.className = 'md-adnotacje__pozycja';
        wiersz.dataset['adnotacja'] = adnotacja.id;
        // Odpowiedź w wątku dostaje wcięcie, żeby wątek dało się przeczytać jako wątek, nie listę zdań.
        if (adnotacja.parentId !== undefined) wiersz.dataset['odpowiedz'] = 'tak';
        wiersz.textContent =
          `${adnotacja.text}` +
          (adnotacja.author === undefined ? '' : ` — ${adnotacja.author}`) +
          (adnotacja.resolved === true ? ' (wątek zamknięty)' : '');
        return wiersz;
      }),
    );
  }

  // Zapis w wątku czyni wskazaną adnotację nadrzędną, a nie celem nadpisania — to dwie różne czynności.
  async function zapisz(wWatku: boolean, zamkniecie: boolean | null): Promise<void> {
    if (bezKompozycji(Command.DesignAnnotationSet)) return;
    if (zamkniecie === null && tresc.kontrolka.value.trim() === '') {
      odpowiedz.pokaz('Pusta uwaga nie mówi niczego, a zajmuje miejsce w wątku.', false);
      return;
    }
    const wskazana = wybor.kontrolka.value;
    if (wWatku && wskazana === '') {
      odpowiedz.pokaz('Wskaż adnotację, na którą odpowiadasz — wątek zaczyna się od uwagi.', false);
      return;
    }
    const zaznaczone = kompozycja.zaznaczone();
    odpowiedz.pokaz('Zapis adnotacji…', true);
    const wynik = await stan.czuwanie.prowadz(
      'zapis adnotacji',
      stan.zrodlo.ustawAdnotacje({
        idKompozycji: kompozycja.idKompozycji(),
        tresc: tresc.kontrolka.value,
        idWarstwy: doWarstwy.kontrolka.checked ? (zaznaczone[0] ?? '') : '',
        idAdnotacji: wWatku ? '' : wskazana,
        idNadrzednej: wWatku ? wskazana : '',
        zamknieta: zamkniecie,
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Zapis adnotacji', wynik.blad), false);
      return;
    }
    const adnotacja = wynik.wynik.annotation;
    zbior = [...zbior.filter((pozycja) => pozycja.id !== adnotacja.id), adnotacja];
    pokazAdnotacje();
    tresc.kontrolka.value = '';
    odpowiedz.pokaz(
      `Rdzeń zapisał adnotację ${adnotacja.id}` +
        (adnotacja.layerId === undefined
          ? ' przypiętą do kompozycji.'
          : ` przypiętą do warstwy ${adnotacja.layerId}.`),
      true,
    );
  }

  /** Zamyka wskazany wątek — zmiana stanu, nie nowa uwaga. */
  async function zamknijWatek(): Promise<void> {
    const wskazana = wybor.kontrolka.value;
    if (wskazana === '') {
      odpowiedz.pokaz('Wskaż adnotację — zamknięcie dotyczy wskazanego wątku.', false);
      return;
    }
    const adnotacja = zbior.find((pozycja) => pozycja.id === wskazana);
    if (adnotacja === undefined) {
      odpowiedz.pokaz('Wskazanej adnotacji nie ma w odczytanym wykazie — odczytaj go ponownie.', false);
      return;
    }
    // Treść idzie ta sama, którą adnotacja niesie — podstawienie z pola okna przepisałoby cudzą uwagę.
    tresc.kontrolka.value = adnotacja.text;
    await zapisz(false, true);
  }

  async function wczytaj(): Promise<void> {
    if (bezKompozycji(Command.DesignAnnotationList)) return;
    odpowiedz.pokaz('Odczyt adnotacji kompozycji…', true);
    const wynik = await stan.czuwanie.prowadz(
      'odczyt adnotacji',
      stan.zrodlo.adnotacje(kompozycja.idKompozycji(), tylkoOtwarte.kontrolka.checked),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt adnotacji', wynik.blad), false);
      return;
    }
    zbior = wynik.wynik.annotations;
    pokazAdnotacje();
    odpowiedz.pokaz(
      zbior.length === 0
        ? 'Ta kompozycja nie ma jeszcze ani jednej adnotacji.'
        : `Adnotacji kompozycji: ${wynik.wynik.total}.`,
      true,
    );
  }

  async function zglosObecnosc(odchodzi: boolean): Promise<void> {
    if (bezKompozycji(Command.DesignPresenceReport)) return;
    const widok = kompozycja.widok();
    const wynik = await stan.czuwanie.prowadz(
      odchodzi ? 'zdjęcie obecności' : 'zgłoszenie obecności',
      stan.zrodlo.zglosObecnosc({
        idKompozycji: kompozycja.idKompozycji(),
        // Położenie pochodzi z przesunięcia kanwy, nie ze śledzenia kursora myszy w pętli.
        x: widok.przesuniecieX,
        y: widok.przesuniecieY,
        zaznaczone: kompozycja.zaznaczone(),
        odchodzi,
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowyBledu(odchodzi ? 'Zdjęcie obecności' : 'Zgłoszenie obecności', wynik.blad),
        false,
      );
      return;
    }
    pokazObecnych(wynik.wynik.participants);
    odpowiedz.pokaz(
      odchodzi
        ? 'Obecność zdjęta — kursor zniknął z kanwy pozostałych.'
        : `Obecność zgłoszona; obecnych na kompozycji: ${wynik.wynik.participants.length}.`,
      true,
    );
  }

  return {
    element,
    wczytaj,
    rozlacz: () => odsubskrybuj(),
  };
}
