import type { AutomationStep } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { poleLiczbowe, poleTekstowe, wiersz } from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import { KLASY_DYMKA, OBJASNIENIA, POZYCJE_BEZ_OBSLUGI } from './etykiety-browser';

/**
 * Limity przebiegu Wykonawcy operującego na stronie — warstwa czwarta
 * Automation Studio. Panel sprawdza scenariusz przed wysłaniem: przekraczający
 * granicę nie wychodzi z okna, a kroki prowadzące poza dozwolone domeny są
 * nazwane po adresie.
 */
export interface LimityPrzebiegu {
  element: HTMLElement;
  /**
   * Powód wstrzymania zapisu; pusty napis znaczy, że scenariusz mieści się
   * w postawionych granicach.
   */
  sprawdz(kroki: readonly AutomationStep[]): string;
}

export function utworzLimityPrzebiegu(pokrycie: PokrycieKomend): LimityPrzebiegu {
  const granica = poleLiczbowe('Górna granica liczby kroków scenariusza', 'bez ograniczenia');
  const wierszGranicy = wiersz('Górna granica liczby kroków scenariusza', granica, {
    klasa: 'dn-pole mb-limity__pole',
    objasnienie: 'Pusta wartość znaczy „bez ograniczenia".',
  });
  const domeny = poleTekstowe({
    etykieta: 'Dozwolone domeny kroków nawigacji',
    podpowiedz: 'adres-strony.pl, konkurent.pl',
    opis: 'Wykaz rozdzielony przecinkami; pusty znaczy „bez ograniczenia".',
  });

  const uwaga = document.createElement('p');
  uwaga.className = 'dn-pole-opis mb-uwaga';
  pokrycie.naOdczyt(() => {
    uwaga.textContent = pokrycie.zdanie(
      POZYCJE_BEZ_OBSLUGI.limityWykonawcy.komenda,
      POZYCJE_BEZ_OBSLUGI.limityWykonawcy.czynnosc,
    );
  });

  const element = document.createElement('section');
  element.className = 'mb-panel mb-limity';
  element.setAttribute('aria-label', 'Limity przebiegu Wykonawcy — warstwa czwarta');
  element.append(
    wierszGranicy,
    domeny.element,
    utworzDymekObjasnienia(OBJASNIENIA.limitKrokow, KLASY_DYMKA),
    uwaga,
  );

  return {
    element,

    sprawdz(kroki) {
      const gorna = Number.parseInt(granica.value, 10);
      if (Number.isFinite(gorna) && gorna > 0 && kroki.length > gorna) {
        return (
          `Scenariusz ma ${kroki.length} kroków przy granicy ${gorna} — zapis wstrzymany ` +
          'w oknie. Podnieś granicę albo skróć scenariusz.'
        );
      }
      const dozwolone = rozbierzDomeny(domeny.kontrolka.value);
      if (dozwolone.length === 0) return '';
      const poza = kroki
        .map((krok) => adresKroku(krok))
        .filter((adres) => adres !== '' && !czyDozwolony(adres, dozwolone));
      if (poza.length === 0) return '';
      return (
        `Kroki prowadzą poza dozwolone domeny: ${poza.join(', ')}. Zapis wstrzymany w oknie — ` +
        'dopisz domenę do wykazu albo usuń krok.'
      );
    },
  };
}

/**
 * Domeny wpisane przez Operatora, bez pustych członów i bez wielkości liter,
 * ponieważ porównanie adresu prowadzi się na postaci znormalizowanej.
 */
function rozbierzDomeny(tresc: string): string[] {
  return tresc
    .split(',')
    .map((czlon) => czlon.trim().toLowerCase())
    .filter((czlon) => czlon !== '');
}

/**
 * Adres, pod który krok prowadzi; pusty napis znaczy „krok nie nawiguje".
 *
 * Treść żądania kroku jest dowolnym obiektem JSON, więc pole adresu sprawdza
 * się po typie. Krok bez adresu nie jest naruszeniem wykazu
 * domen — jest krokiem innego rodzaju.
 */
function adresKroku(krok: AutomationStep): string {
  const tresc = krok.params;
  if (typeof tresc !== 'object' || tresc === null || Array.isArray(tresc)) return '';
  const adres = (tresc as Record<string, unknown>)['url'];
  return typeof adres === 'string' ? adres.trim() : '';
}

/**
 * Czy adres należy do którejś z dozwolonych domen albo jej poddomeny; pusty
 * wykaz domen niczego nie ogranicza.
 */
function czyDozwolony(adres: string, dozwolone: readonly string[]): boolean {
  const nazwa = nazwaHosta(adres);
  if (nazwa === '') return false;
  return dozwolone.some((domena) => nazwa === domena || nazwa.endsWith(`.${domena}`));
}

/**
 * Nazwa hosta adresu; pusta, gdy treść adresem nie jest, ponieważ krok bez
 * adresu nie podlega sprawdzeniu wykazu domen.
 */
function nazwaHosta(adres: string): string {
  try {
    return new URL(adres).hostname.toLowerCase();
  } catch {
    return '';
  }
}
