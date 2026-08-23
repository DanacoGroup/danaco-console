import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk } from '../../modele/kontrolki-formularza';
import { odczytajBase64 } from './tresc-base64';
import { KOD_MODULU } from './zrodlo-otoczenia';
import type { StanBiblioteki } from './stan-biblioteki';

/**
 * Wgranie pliku do repozytorium wiedzy — `library.file.upload`.
 *
 * Treść idzie base64, bo takie pole niesie kontrakt (`contentBase64`), a plik
 * wskazany z urządzenia nie ma ścieżki widocznej dla rdzenia: przeglądarka
 * podaje wyłącznie nazwę. Pole `sourcePath` zostaje więc puste — wypełnienie go
 * nazwą pliku byłoby danymi zmyślonymi.
 *
 * Odczyt pliku dzieje się przed wywołaniem i może się nie udać (plik zniknął,
 * odczyt odrzucony). Niepowodzenie odczytu wraca tą samą drogą co odmowa
 * rdzenia — jednym zdaniem w wierszu odpowiedzi, bez wyjątku wywracającego widok.
 *
 * Powodzenie komendy nie jest powodzeniem wgrania: rdzeń bez magazynu treści
 * zakłada wpis o pliku i odmawia jego podglądu. Okno pyta więc rdzeń o treść
 * zaraz po wgraniu (`dostepnosc-tresci.ts`) i podaje to, co rdzeń odpowiedział —
 * kosztem jednej dodatkowej komendy, bo kontrakt nie ma tańszego świadka.
 */
export interface WgraniePliku {
  /** Kontrolki osadzane w pasku okna. */
  element: HTMLElement;
  /** Otwiera wybór pliku — używa tego również przycisk stanu pustego. */
  wskazPlik(): void;
}

export function utworzWgraniePliku(
  stan: StanBiblioteki,
  odpowiedz: (tresc: string, powodzenie: boolean) => void,
): WgraniePliku {
  const wybor = document.createElement('input');
  wybor.type = 'file';
  wybor.hidden = true;

  const dodaj = przycisk('Dodaj plik', 'dn-btn dn-btn--sm dn-btn--atrament');
  dodaj.dataset['czynnosc'] = 'wgraj';
  dodaj.addEventListener('click', () => wybor.click());

  wybor.addEventListener('change', () => {
    const plik = wybor.files?.[0];
    if (plik === undefined) return;
    // Kontrolka wyboru pliku zgłasza `change` tylko przy zmianie wartości.
    // Wyczyszczenie jej po przejęciu uchwytu pozwala wskazać ten sam plik
    // ponownie, na przykład po nieudanym wgraniu.
    wybor.value = '';
    void wgraj(plik);
  });

  async function wgraj(plik: File): Promise<void> {
    odpowiedz(`Wgrywanie pliku „${plik.name}"…`, true);
    const tresc = await odczytajBase64(plik);
    if (tresc === null) {
      odpowiedz(`Odczyt pliku „${plik.name}" nie powiódł się — nic nie zostało wysłane.`, false);
      return;
    }
    const wynik = await stan.zrodlo.wgraj({
      name: plik.name,
      sourceModuleId: KOD_MODULU,
      contentBase64: tresc,
      ...(plik.type === '' ? {} : { mimeType: plik.type }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz(opisOdmowy('Wgranie pliku', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const wgrany = wynik.wynik.file;
    stan.wchlon(wgrany);
    odpowiedz(`Wpis pliku „${wgrany.name}" założony — pytam rdzeń o jego treść…`, true);
    const stanTresci = await stan.zbadajTresc(wgrany.id);
    if (stanTresci.werdykt === 'brak') {
      odpowiedz(
        `Plik „${wgrany.name}" jest w wykazie repozytorium, ale jego TREŚĆ do niego nie ` +
          `trafiła — ${stanTresci.powod}. Wgranie nie jest domknięte: podgląd, wersje ` +
          'i przeniesienie do modułu nie mają czego pokazać.',
        false,
      );
      return;
    }
    // Odpowiedź ze wskazaniem miejsca treści nie jest odmową: taki podgląd
    // wraca zarówno dla treści leżącej w magazynie, jak i dla wskazania
    // w próżnię (`dostepnosc-tresci.ts`). Zdanie mówi więc, co przyszło,
    // i nie orzeka ani wgrania, ani jego braku.
    if (stanTresci.werdykt === 'odwolanie') {
      odpowiedz(
        `Wpis pliku „${wgrany.name}" stoi w wykazie, a rdzeń na pytanie o treść oddał ` +
          `wskazanie zamiast bajtów — ${stanTresci.powod}. Wgranie zostaje ` +
          'nierozstrzygnięte: podgląd nie pokazał treści tego pliku.',
        false,
      );
      return;
    }
    // Powodzenie mówi się wyłącznie po werdykcie `osiagalna`. Werdykt `odmowa`
    // znaczy, że rdzeń treści nie oddał i o niej nie orzekł — okno nie orzeka
    // wtedy ani obecności treści, ani jej braku.
    if (stanTresci.werdykt !== 'osiagalna') {
      odpowiedz(
        `Wpis pliku „${wgrany.name}" stoi w wykazie, ale rdzeń nie odpowiedział, czy ma ` +
          `jego treść — ${stanTresci.powod}. Wgranie zostaje nierozstrzygnięte: spróbuj ` +
          'podglądu ponownie, zanim uznasz plik za wgrany.',
        false,
      );
      return;
    }
    odpowiedz(`Plik „${wgrany.name}" wgrany — rdzeń oddaje jego treść do podglądu.`, true);
  }

  const element = document.createElement('span');
  element.className = 'ml-wgranie';
  element.append(dodaj, wybor);

  return { element, wskazPlik: () => wybor.click() };
}
