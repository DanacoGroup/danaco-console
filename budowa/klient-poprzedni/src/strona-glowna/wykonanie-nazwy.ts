import type { Kanal } from '../protokol/kanal';
import { zadajKopieSesji, zadajZmianeNazwySesji } from '../protokol/nazwa-sesji';
import { zadajPrzypisanieProjektu } from '../protokol/projekt-sesji';
import type { CzynnosciSesji, MeldunekCzynnosci } from './czynnosci-sesji';
import { odmowa, nazwa } from './meldunki-sesji';
import { zapytajONazwe } from './pytanie-o-nazwe';

/**
 * Trzy czynności historii, które przed wysłaniem komendy pytają o napis.
 *
 * Jedna odpowiedzialność: zebranie nazwy od Operatora i dopiero z nią wyjście
 * na kontrakt. Stoją osobno od pozostałych czynności, bo mają wspólny kształt
 * — pytanie, odmowa, komenda, meldunek — i wspólny warunek zejścia: zamknięte
 * pytanie kończy rzecz bez żadnego żądania.
 *
 * Zamknięcie pytania zwraca `null`, nie napis pusty. Odmowa Operatora nie jest
 * odmową rdzenia i nie ma o niej czego meldować, a napis pusty bywa odpowiedzią
 * sensowną — przy kopii oddaje nadanie nazwy rdzeniowi — więc te dwa przypadki
 * nie mogą się zlać.
 */
export function czynnosciZNazwa(kanal: Kanal, odswiez: () => void): CzynnosciSesji {
  return {
    async zmienNazwe(wpis): Promise<MeldunekCzynnosci> {
      const tytul = await zapytajONazwe({
        tytul: 'Zmiana nazwy sesji',
        etykieta: 'Nowa nazwa',
        wartosc: wpis.sesja.title ?? '',
        opis: 'Zmienia wyłącznie nazwę w historii; zapis rozmowy zostaje nietknięty.',
        napisZatwierdzenia: 'Zmień nazwę',
      });
      if (tytul === null) return null;
      if (tytul === '') return 'Nazwa sesji nie może być pusta.';

      const wynik = await zadajZmianeNazwySesji(kanal, { sessionId: wpis.sesja.id, title: tytul });
      if (!wynik.udany) return odmowa(wynik, 'zmiana nazwy sesji');
      odswiez();
      return `Sesja nosi teraz nazwę „${wynik.wynik?.session.title ?? tytul}”.`;
    },

    async skopiuj(wpis): Promise<MeldunekCzynnosci> {
      const tytul = await zapytajONazwe({
        tytul: 'Kopia sesji',
        etykieta: 'Nazwa kopii',
        opis: 'Pole puste zostawia nadanie nazwy rdzeniowi — weźmie nazwę źródła z dopiskiem.',
        napisZatwierdzenia: 'Kopiuj',
      });
      if (tytul === null) return null;

      const wynik = await zadajKopieSesji(kanal, {
        sessionId: wpis.sesja.id,
        ...(tytul === '' ? {} : { title: tytul }),
      });
      if (!wynik.udany) return odmowa(wynik, 'kopiowanie sesji');
      odswiez();
      const tresc = wynik.wynik;
      const kopia = tresc?.session.title ?? tresc?.session.id ?? '';
      return `Powstała kopia „${kopia}” — przeniesiono ${tresc?.copiedMessages ?? 0} wiadomości.`;
    },

    async przypiszProjekt(wpis): Promise<MeldunekCzynnosci> {
      const projekt = await zapytajONazwe({
        tytul: 'Sesja do projektu',
        etykieta: 'Nazwa projektu',
        opis: 'Projekt o tej nazwie zostanie założony, jeśli jeszcze nie istnieje.',
        napisZatwierdzenia: 'Przenieś',
      });
      if (projekt === null) return null;
      if (projekt === '') return 'Nazwa projektu nie może być pusta.';

      const wynik = await zadajPrzypisanieProjektu(kanal, {
        sessionIds: [wpis.sesja.id],
        projectName: projekt,
      });
      if (!wynik.udany) return odmowa(wynik, 'przeniesienie sesji do projektu');
      odswiez();
      return (wynik.wynik?.movedIds.length ?? 0) > 0
        ? `Sesja „${nazwa(wpis)}” trafiła do projektu „${projekt}”.`
        : 'Rdzeń nie przeniósł żadnej sesji.';
    },
  };
}
