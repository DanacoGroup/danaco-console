import { Command, type DesignAsset } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { przycisk, ustawPozycje, poleWyboru, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { BRAKI } from './etykiety-designu';
import { przyciskBrakuDrogi } from './brak-drogi';
import { nazwaZasobu, opisZasobu } from './karta-zasobu';
import { skutekPrzekazania } from './skutek-designu';
import type { StanDesignu } from './stan-designu';

/**
 * Panel metadanych jednego zasobu wskazanego w Assets Panel wraz
 * z przekazaniem go do modułu docelowego.
 *
 * „Wyślij do modułu docelowego" idzie komendą `context.transfer` — jedyną
 * zbudowaną w kontrakcie drogą przeniesienia kompletu kontekstu. Wykaz modułów
 * docelowych przychodzi z rdzenia (`module.list`), nie z kopii katalogu po
 * stronie klienta.
 *
 * `ContextBundle` nie ma pola zasobu wizualnego, więc zasób jedzie polem
 * `documentIds` wraz z poleceniem nazywającym go po imieniu — to przybliżenie
 * kontraktu, nie jego pełne pokrycie.
 *
 * Potwierdzenie opisuje okno, które wróciło, a nie moduł zamówiony w polu
 * wyboru: `context.transfer` nie sprawdza katalogu modułów — przepisuje
 * `targetModuleId` do okna docelowego i oddaje `transferred: true` także dla
 * kodu, którego katalog nie zna. Jedynym polem mówiącym, gdzie zasób wylądował,
 * jest `window.moduleId` odpowiedzi; rozbieżność wobec zamówienia jest odmową
 * (`skutek-designu.ts`).
 */
export interface PanelMetadanych {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzPanelMetadanych(stan: StanDesignu): PanelMetadanych {
  const tytul = document.createElement('h4');
  tytul.className = 'md-panel__tytul';
  tytul.textContent = 'Metadane zasobu';

  const opis = document.createElement('p');
  opis.className = 'dn-pole-opis md-panel__opis';

  const warianty = document.createElement('ul');
  warianty.className = 'md-warianty';

  const cel = poleWyboru({
    etykieta: 'Moduł docelowy przekazania',
    opis: 'Wykaz z katalogu rdzenia (module.list). Przekazanie idzie komendą context.transfer.',
  }, []);

  const wyslij = przycisk('Wyślij do modułu docelowego', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedz = utworzWierszOdpowiedzi();

  const rzadBrakow = document.createElement('div');
  rzadBrakow.className = 'md-braki__rzad';
  // Eksport zbiorczy ma już komendę w kontrakcie, ale nie ma uchwytu w rdzeniu,
  // więc stoi tu jako kontrolka nazywająca stan drogi. Etykietowanie, wgranie
  // i usunięcie zasobu mają własne kontrolki wykonujące
  // (`nadanie-etykiet.ts`, `wgranie-zasobu.ts`, `czynnosci-zasobu.ts`).
  rzadBrakow.append(przyciskBrakuDrogi(BRAKI.eksportZbiorczy));

  const element = document.createElement('section');
  element.className = 'md-panel md-metadane';
  element.append(tytul, opis, warianty, cel.element, wyslij, odpowiedz.element, rzadBrakow);

  async function przekaz(): Promise<void> {
    const zasob = stan.wybrany();
    if (zasob === null) {
      odpowiedz.pokaz('Wskaż zasób w wykazie — przekazanie dotyczy jednego zasobu.', false);
      return;
    }
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(`Komenda ${Command.ContextTransfer} wymaga okna źródłowego. ${stan.opisOkna()}`, false);
      return;
    }
    if (cel.kontrolka.value === '') {
      odpowiedz.pokaz('Wskaż moduł docelowy — katalog modułów przyszedł pusty albo odmową.', false);
      return;
    }
    odpowiedz.pokaz(`Przekazanie zasobu „${nazwaZasobu(zasob)}"…`, true);
    // Pod czuwaniem — po zerwaniu gniazda wiersz mówi prawdę o braku
    // rozstrzygnięcia zamiast stać na zdaniu o przekazaniu w toku.
    const wynik = await stan.czuwanie.prowadz(
      'przekazanie zasobu',
      stan.zaplecze.przekaz({
        idOknaZrodlowego: stan.idOkna(),
        kodModuluDocelowego: cel.kontrolka.value,
        komplet: {
          prompt: `Zasób wizualny modułu Design: ${nazwaZasobu(zasob)} (${zasob.kind}).`,
          documentIds: [zasob.id],
        },
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
      odpowiedz.pokaz(opisOdmowyBledu('Przekazanie zasobu', wynik.blad), false);
      return;
    }
    const skutek = skutekPrzekazania(cel.kontrolka.value, wynik.wynik, nazwaZasobu(zasob));
    odpowiedz.pokaz(skutek.zdanie, skutek.udany);
  }

  wyslij.addEventListener('click', () => void przekaz());

  return {
    element,

    odswiez() {
      ustawPozycje(
        cel.kontrolka,
        stan.celePrzekazania().map((modul) => ({ wartosc: modul.code, etykieta: modul.name })),
      );
      const zasob = stan.wybrany();
      if (zasob === null) {
        opis.textContent = 'Żaden zasób nie jest wskazany — metadanych nie ma czego dotyczyć.';
        warianty.replaceChildren();
        return;
      }
      opis.textContent = `${nazwaZasobu(zasob)} — ${opisZasobu(zasob)}`;
      warianty.replaceChildren(...wierszeWariantow(stan.zasoby(), zasob));
    },
  };
}

/** Historia wariantów: zasoby wskazujące ten sam pierwowzór albo ten zasób. */
function wierszeWariantow(
  wszystkie: readonly DesignAsset[],
  zasob: DesignAsset,
): readonly HTMLElement[] {
  const pierwowzor = zasob.variantOfAssetId ?? zasob.id;
  const rodzina = wszystkie.filter(
    (wpis) => wpis.id === pierwowzor || wpis.variantOfAssetId === pierwowzor,
  );
  if (rodzina.length < 2) return [wiersz('Zasób nie ma jeszcze wariantów w wykazie wczytanym.')];
  return rodzina.map((wpis) =>
    wiersz(`${nazwaZasobu(wpis)}${wpis.id === zasob.id ? ' — wskazany' : ''}`),
  );
}

function wiersz(tresc: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'md-warianty__wiersz';
  element.textContent = tresc;
  return element;
}
