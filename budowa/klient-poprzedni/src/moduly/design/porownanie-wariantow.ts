import type { DesignAsset } from '../../../../shared/contract';
import { poleWyboru, ustawPozycje } from '../../modele/kontrolki-formularza';
import { nazwaZasobu } from './karta-zasobu';
import { czyRdzenZnaObraz, slowoModelu } from './slowo-modelu';

/**
 * Porównanie „przed / po" dwóch zasobów jednej rodziny wariantów zestawia pola
 * opisowe obu stron. Rodzinę wyznacza pole `DesignAsset.variantOfAssetId`:
 * pierwowzór wraz ze wszystkimi zasobami, które go wskazują.
 */
export interface PorownanieWariantow {
  element: HTMLElement;
  /** Nanosi rodzinę zasobu wskazanego; `null` gasi porównanie. */
  odswiez(zasoby: readonly DesignAsset[], wskazany: DesignAsset | null): void;
}

export function utworzPorownanieWariantow(): PorownanieWariantow {
  const wybor = poleWyboru(
    {
      etykieta: 'Wariant po lewej („przed")',
      opis:
        'Zasoby wspólnego pierwowzoru — pole variantOfAssetId. Po prawej stoi zasób ' +
        'wskazany w Assets Panel.',
    },
    [],
  );

  const lewa = kolumna('przed');
  const prawa = kolumna('po');

  const para = document.createElement('div');
  para.className = 'md-porownanie__para';
  para.append(lewa.element, prawa.element);

  const zdanie = document.createElement('p');
  zdanie.className = 'dn-pole-opis md-porownanie__zdanie';

  const element = document.createElement('section');
  element.className = 'md-porownanie';
  element.setAttribute('aria-label', 'Porównanie wariantów zasobu');
  element.append(wybor.element, para, zdanie);

  /** Rodzina zapamiętana między odświeżeniami — wybór ma z czego czytać. */
  let rodzina: readonly DesignAsset[] = [];

  function nanies(wskazany: DesignAsset | null): void {
    const przed = rodzina.find((wpis) => wpis.id === wybor.kontrolka.value) ?? null;
    lewa.pokaz(przed);
    prawa.pokaz(wskazany);
  }

  wybor.kontrolka.addEventListener('change', () => nanies(prawa.zasob()));

  return {
    element,

    odswiez(zasoby, wskazany) {
      if (wskazany === null) {
        rodzina = [];
        ustawPozycje(wybor.kontrolka, []);
        zdanie.textContent = 'Porównanie dotyczy zasobu wskazanego — bez wskazania nie ma dwóch stron.';
        nanies(null);
        return;
      }
      const pierwowzor = wskazany.variantOfAssetId ?? wskazany.id;
      rodzina = zasoby.filter(
        (wpis) => wpis.id === pierwowzor || wpis.variantOfAssetId === pierwowzor,
      );
      const inne = rodzina.filter((wpis) => wpis.id !== wskazany.id);
      const poprzedni = wybor.kontrolka.value;
      ustawPozycje(
        wybor.kontrolka,
        inne.map((wpis) => ({ wartosc: wpis.id, etykieta: nazwaZasobu(wpis) })),
      );
      if (inne.some((wpis) => wpis.id === poprzedni)) wybor.kontrolka.value = poprzedni;

      zdanie.textContent =
        inne.length === 0
          ? `Zasób „${nazwaZasobu(wskazany)}" nie ma w wykazie wczytanym ani jednego rodzeństwa — ` +
            'porównywać nie ma z czym.'
          : 'Porównywane są SŁOWA modelu i pola opisowe, nie obrazy. Rdzeń obrazy generuje ' +
            'i bajty ma u siebie, a droga po nie jest już w kontrakcie — brakuje jej uchwytu ' +
            'w rdzeniu, więc obu stron nie ma dziś czym wczytać ani nałożyć suwakiem.';
      nanies(wskazany);
    },
  };
}

/**
 * Buduje jedną stronę porównania: tytuł strony, nazwę zasobu, słowo modelu oraz
 * stan treści obrazu. Strona przechowuje własny zasób i oddaje go na żądanie,
 * a wartość pusta gasi treść kolumny.
 */
function kolumna(strona: string): {
  element: HTMLElement;
  pokaz(zasob: DesignAsset | null): void;
  zasob(): DesignAsset | null;
} {
  const tytul = document.createElement('h5');
  tytul.className = 'md-porownanie__tytul';

  const tresc = document.createElement('p');
  tresc.className = 'md-porownanie__tresc';

  const obraz = document.createElement('p');
  obraz.className = 'dn-pole-opis md-porownanie__obraz';

  const element = document.createElement('div');
  element.className = 'md-porownanie__strona';
  element.dataset['strona'] = strona;
  element.append(tytul, tresc, obraz);

  let biezacy: DesignAsset | null = null;

  return {
    element,
    zasob: () => biezacy,

    pokaz(zasob) {
      biezacy = zasob;
      if (zasob === null) {
        tytul.textContent = `${strona} — brak`;
        tresc.textContent = '';
        obraz.textContent = '';
        return;
      }
      const slowo = slowoModelu(zasob);
      tytul.textContent = `${strona} — ${nazwaZasobu(zasob)}`;
      tresc.textContent = slowo === '' ? 'Model nie oddał opisu słownego tego zasobu.' : slowo;
      obraz.textContent = czyRdzenZnaObraz(zasob)
        ? `Treść obrazu: ${zasob.uri ?? ''}`
        : 'Treść obrazu: rdzeń jej nie zna.';
    },
  };
}
