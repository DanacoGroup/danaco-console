import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk } from '../../modele/kontrolki-formularza';
import { odczytajBase64 } from './tresc-base64';
import { KOD_MODULU } from './zrodlo-otoczenia';
import type { StanBiblioteki } from './stan-biblioteki';

/**
 * Wgranie pliku do repozytorium wiedzy komendą `library.file.upload`. Treść pliku
 * jedzie w zapisie base64, a po założeniu wpisu okno pyta rdzeń o dostępność
 * treści i dopiero jego odpowiedź rozstrzyga o powodzeniu wgrania.
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
    // Wyczyszczenie kontrolki pozwala wskazać ten sam plik ponownie.
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
    // Wskazanie miejsca treści nie jest ani odmową, ani dowodem wgrania.
    if (stanTresci.werdykt === 'odwolanie') {
      odpowiedz(
        `Wpis pliku „${wgrany.name}" stoi w wykazie, a rdzeń na pytanie o treść oddał ` +
          `wskazanie zamiast bajtów — ${stanTresci.powod}. Wgranie zostaje ` +
          'nierozstrzygnięte: podgląd nie pokazał treści tego pliku.',
        false,
      );
      return;
    }
    // Powodzenie okno orzeka wyłącznie po werdykcie osiągalności treści.
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
