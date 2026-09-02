// Kolekcje zasobów okna: wykaz, założenie i przypisanie wskazanego zasobu.
// Droga przeniesiona z poprzedniego klienta, gdzie stała w panelu kolekcji.

import { Command } from '../../../shared/contract.ts';
import { oglos } from './ogloszenie.ts';
import { chwila, poproszony, wykonany, type Kontekst, type Wiersz } from './design-wspolne.ts';
import { dopiszDoKorzenia, pokazWykaz } from './design-wykaz.ts';

export function zwiazKolekcje(kontekst: Kontekst): void {
  kontekst.stan.odswiezenia.set('kolekcje', () => {
    void pokazKolekcje(kontekst);
  });
  dopiszDoKorzenia({
    nazwa: 'Kolekcje zasobów',
    opis: 'wykaz, założenie, przypisanie',
    otworz: (biezacy) => {
      void pokazKolekcje(biezacy);
    },
  });
}

export async function pokazKolekcje(kontekst: Kontekst): Promise<void> {
  if (kontekst.stan.idOkna === '') return;
  // Zawężenia do zasobu nie ma: do każdej kolekcji okna wolno go dołożyć.
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignCollectionList, {
    windowId: kontekst.stan.idOkna,
  });
  const wiersze: Wiersz[] = (odpowiedz?.collections ?? []).map((kolekcja) => ({
    tekst: kolekcja.name,
    meta: `${String(kolekcja.assetCount)} zasobów · ${chwila(kolekcja.updatedAt)}`,
    plakietka: kolekcja.description ?? '',
    kropka: 'neutralna',
    naKlik: () => {
      void przypisz(kontekst, kolekcja.id);
    },
  }));
  wiersze.unshift({
    tekst: 'Załóż kolekcję z wskazanego zasobu',
    kropka: 'sygnal',
    naKlik: () => {
      void zaloz(kontekst);
    },
  });
  pokazWykaz(kontekst, 'Kolekcje zasobów', wiersze, 'Okno nie ma jeszcze kolekcji.');
}

async function zaloz(kontekst: Kontekst): Promise<void> {
  const odpowiedz = await wykonany(kontekst.kanal, Command.DesignCollectionCreate, {
    windowId: kontekst.stan.idOkna,
    name: `Kolekcja ${chwila(Date.now())}`,
  }, 'Założenie kolekcji');
  if (odpowiedz === null) return;
  if (kontekst.stan.idZasobu !== '') await przypisz(kontekst, odpowiedz.collection.id);
  else await pokazKolekcje(kontekst);
}

async function przypisz(kontekst: Kontekst, idKolekcji: string): Promise<void> {
  if (kontekst.stan.idZasobu === '') {
    oglos('Design', 'Wskaż zasób w Assets Panel, zanim trafi do kolekcji.', 'ostrzezenie');
    return;
  }
  await wykonany(kontekst.kanal, Command.DesignCollectionAssign, {
    collectionId: idKolekcji,
    assetIds: [kontekst.stan.idZasobu],
  }, 'Przypisanie do kolekcji');
  await pokazKolekcje(kontekst);
}
