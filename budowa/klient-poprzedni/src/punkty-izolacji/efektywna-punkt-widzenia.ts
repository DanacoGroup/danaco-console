import type { IsolationPolicyPreviewRequest } from '../../../shared/contract';
import { pole, wiersz } from '../modele/kontrolki-formularza-braki';
import type { Kanal } from '../protokol/kanal';
import { idSesji } from './komenda';
import { NAZWY_WARSTW, type StanWarstwy } from './stan-warstwy';
import type { StanZasiegu } from './stan-zasiegu';

/** Punkt widzenia podglądu polityki — czego dotyczy nazwa „polityka efektywna" w oknie punktów izolacji. */
export interface PunktWidzeniaPodgladu {
  /** Panel montowany nad tabelą podglądu. */
  element: HTMLElement;
  /** Treść żądania `isolation.policy.preview` złożona z wyboru Operatora. */
  zadanie(): IsolationPolicyPreviewRequest;
  /** Zdanie opisujące, czego dotyczy odczyt — do wstępu nad tabelą. */
  opis(): string;
}

export function utworzPunktWidzenia(
  kanal: Kanal,
  warstwa: StanWarstwy,
  zasieg: StanZasiegu,
  naZmiane: () => void,
): PunktWidzeniaPodgladu {
  const oknoKomunikacji = pole('Okno komunikacji', 'puste — bez zawężenia do okna');

  const odswiez = document.createElement('button');
  odswiez.type = 'button';
  odswiez.className = 'dn-btn';
  odswiez.textContent = 'Odczytaj podgląd dla tego punktu widzenia';
  odswiez.addEventListener('click', () => naZmiane());

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'pi-widzenie__opis';
  wyjasnienie.textContent =
    'Podgląd bez podanego punktu widzenia schodzi w rdzeniu do poziomu globalnego — pokazywałby ' +
    'politykę całej platformy pod nazwą „efektywna”. Kartę sesji i warstwę okno podaje samo, poziom ' +
    'zasięgu i byt bierze z selektora zasięgu w panelu lewym; okno komunikacji wskazuje Operator, ' +
    'bo klient go nie zgaduje.';

  const stanOdczytu = document.createElement('p');
  stanOdczytu.className = 'pi-widzenie__opis';

  const element = document.createElement('div');
  element.className = 'pi-widzenie';
  element.append(
    wyjasnienie,
    wiersz('Okno komunikacji (identyfikator)', oknoKomunikacji, { klasa: 'dn-pole' }),
    odswiez,
    stanOdczytu,
  );

  function zadanie(): IsolationPolicyPreviewRequest {
    const byt = zasieg.bytDoZadania();
    const sessionId = idSesji(kanal);
    return {
      scope: zasieg.zasieg(),
      ...(byt === undefined ? {} : { scopeId: byt }),
      ...(sessionId === undefined ? {} : { sessionId }),
      ...(oknoKomunikacji.value === '' ? {} : { windowId: oknoKomunikacji.value }),
      layer: warstwa.warstwa(),
    };
  }

  function opis(): string {
    const sessionId = idSesji(kanal);
    const czesci = [
      zasieg.opis(),
      sessionId === undefined ? 'bez karty sesji (rdzeń jej jeszcze nie założył)' : `karta sesji ${sessionId}`,
      oknoKomunikacji.value === '' ? 'bez zawężenia do okna' : `okno ${oknoKomunikacji.value}`,
      `warstwa „${NAZWY_WARSTW[warstwa.warstwa()]}”`,
    ];
    return `Odczyt dotyczy: ${czesci.join(', ')}.`;
  }

  oknoKomunikacji.addEventListener('change', () => {
    stanOdczytu.textContent = `${opis()} Naciśnij „Odczytaj podgląd”, żeby zapytać rdzeń o ten punkt widzenia.`;
  });

  stanOdczytu.textContent = opis();

  return { element, zadanie, opis };
}
