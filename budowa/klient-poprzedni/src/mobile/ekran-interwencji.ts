import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';
import { utworzArkuszDrog, type ArkuszDrog } from './arkusz-drog';
import { utworzMagazynKwitow, type MagazynKwitow } from './kwit-decyzji';
import { utworzPasekKwitu, type PasekKwitu } from './pasek-kwitu';
import { kolejkaDecyzji } from './port-kolejki-decyzji';
import { skrocObraz, zlozPozycjeDecyzji, type PozycjaDecyzji } from './pozycje-decyzji';
import { utworzWywolaniaInterwencji } from './wywolania-interwencji';
import { utworzZrodloInterwencji, type ZrodloInterwencji } from './zrodlo-interwencji';

/**
 * Ekran interwencji — telefon jako kanał pojedynczej decyzji, nie mniejszy
 * pulpit: jedna kolumna, karta pozycji jako całość celem dotknięcia, arkusz
 * dróg u dołu. Odmowa jednego odczytu nie gasi ekranu — wykaz stoi dalej na
 * odczytach, które doszły.
 */
export interface EkranInterwencji {
  element: HTMLElement;
  /** Pyta rdzeń i przerysowuje wykaz. */
  odswiez(): Promise<void>;
  /** Zdejmuje nasłuch zdarzeń — wołane przy rozłączeniu kanału. */
  rozlacz(): void;
  /** Pozycje z ostatniego odczytu — do sprawdzianów widoku. */
  pozycje(): readonly PozycjaDecyzji[];
}

export interface OpisEkranu {
  kanal: Kanal;
  /** Magazyn kwitów; wstrzykiwany dla sprawdzianu, domyślnie pamięć przeglądarki. */
  magazyn?: MagazynKwitow;
}

export function utworzEkranInterwencji(opis: OpisEkranu): EkranInterwencji {
  const zrodlo: ZrodloInterwencji = utworzZrodloInterwencji(opis.kanal);
  const wywolania = utworzWywolaniaInterwencji(opis.kanal);
  const magazyn = opis.magazyn ?? utworzMagazynKwitow();

  let ostatnie: PozycjaDecyzji[] = [];
  let wTrakcie = false;

  const stan = document.createElement('p');
  stan.className = 'mb-ekran__stan';
  stan.setAttribute('aria-live', 'polite');

  const naglowek = document.createElement('ul');
  naglowek.className = 'mb-ekran__naglowek';

  const zrodloKolejki = document.createElement('p');
  zrodloKolejki.className = 'mb-ekran__zrodlo';

  const wykaz = document.createElement('div');
  wykaz.className = 'mb-ekran__wykaz';

  const pasek: PasekKwitu = utworzPasekKwitu(magazyn);
  const arkusz: ArkuszDrog = utworzArkuszDrog({
    wywolania,
    poKwicie: (kwit) => {
      pasek.dopisz(kwit);
      // Po każdej drodze pytamy rdzeń ponownie: wykaz ma pokazywać stan PO
      // decyzji, a nie ten sprzed niej.
      void odswiez();
    },
  });

  const element = document.createElement('div');
  element.className = 'mb-ekran';
  element.append(stan, naglowek, zrodloKolejki, wykaz, arkusz.element, pasek.element);

  /** Karta pozycji — celem dotknięcia jest cała karta, nie osobny przycisk. */
  function kartaPozycji(pozycja: PozycjaDecyzji): HTMLElement {
    const karta = document.createElement('button');
    karta.type = 'button';
    karta.className = 'dn-karta dn-karta--klikalna mb-karta';
    karta.dataset['rodzaj'] = pozycja.rodzaj;

    const tytul = document.createElement('span');
    tytul.className = 'mb-karta__tytul';
    tytul.textContent = pozycja.naglowek;

    const pierwszy = document.createElement('span');
    pierwszy.className = 'mb-karta__kontekst';
    pierwszy.textContent = pozycja.kontekst[0] ?? `Źródło: ${pozycja.zrodlo}.`;

    karta.append(tytul, pierwszy);
    karta.addEventListener('click', () => arkusz.pokaz(pozycja));
    return karta;
  }

  async function odswiez(): Promise<void> {
    if (wTrakcie) return;
    wTrakcie = true;
    stan.dataset['stan'] = 'ladowanie';
    stan.textContent = 'Pytam rdzeń, co czeka na decyzję…';

    try {
      const obraz = await zrodlo.zbierzObraz();
      const skrot = skrocObraz(obraz);
      naglowek.replaceChildren(
        ...skrot.zdania.map((zdanie) => {
          const wiersz = document.createElement('li');
          wiersz.textContent = zdanie;
          return wiersz;
        }),
      );

      const zPortu = await kolejkaDecyzji().odczytaj();
      zrodloKolejki.textContent = zPortu.dostepna
        ? `Kolejka eskalacji: ${kolejkaDecyzji().nazwa}.`
        : (zPortu.powodBraku ?? '');

      ostatnie = [...zPortu.pozycje, ...zlozPozycjeDecyzji(obraz, Date.now())];
      wykaz.replaceChildren(...ostatnie.map(kartaPozycji));

      if (ostatnie.length > 0) {
        stan.dataset['stan'] = 'tresc';
        stan.textContent = `Czeka na ciebie: ${ostatnie.length}.`;
      } else if (skrot.odmowy > 0) {
        // Odmowy były, więc „nic nie czeka" byłoby zdaniem nieprawdziwym.
        stan.dataset['stan'] = 'blad';
        stan.textContent =
          `Nie wiem, czy coś czeka: rdzeń odmówił ${skrot.odmowy} z odczytów. ` +
          opisOdmowyBledu('Powód pierwszego odczytu', obraz.procesy.blad);
      } else {
        stan.dataset['stan'] = 'pusto';
        stan.textContent = 'Nic nie czeka na twoją decyzję. Rdzeń odpowiedział na każdy odczyt.';
      }

      pasek.odswiez();
    } finally {
      wTrakcie = false;
    }
  }

  // Zdarzenia rdzenia przerysowują wykaz same; nasłuch stoi na trzech zdarzeniach, które go mają.
  const odsubskrybuj = [
    zrodlo.naPostep(() => void odswiez()),
    zrodlo.naKolejke(() => void odswiez()),
    zrodlo.naStanOkna(() => void odswiez()),
  ];

  return {
    element,
    odswiez,
    rozlacz() {
      for (const zdejmij of odsubskrybuj) zdejmij();
    },
    pozycje: () => ostatnie,
  };
}
