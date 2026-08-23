import { elementGodla } from '../ikony/ikony';

/**
 * Godło marki i tożsamość projektu pokazywane na pasku aplikacji.
 *
 * Godło buduje warstwa znaku marki (`ikony/ikony.ts` → `ikony/marka.ts`).
 * Znak ma barwy własne, więc nie dziedziczy barwy tekstu i nie zmienia się
 * z motywem; odmianę dobiera podłoże, a pasek aplikacji jest atramentowy
 * w obu motywach — stąd podłoże ciemne. Droga przez `elementGodla` niesie
 * ponadto próg odmiany uproszczonej: od 16 px w dół znak przechodzi na jeden grot.
 *
 * Nazwa produktu jest stała, nazwa projektu przychodzi z opisu okna, czyli
 * z konfiguracji budowania (`okno-komunikacji/opis-okna.ts`). Gdy nazwa
 * projektu równa się nazwie produktu, na pasek nie wchodzi wcale —
 * powtórzenie tej samej nazwy nic nie wnosi, a zajmuje miejsce.
 *
 * Blok nie jest kontrolką i nią nie udaje — nie prowadzi nigdzie i nie ma
 * kursora wskazującego.
 */
/** Nazwa produktu — jedno miejsce, w którym pada na pasku. */
const NAZWA_PRODUKTU = 'Danaco Console';

export function utworzGodlo(projekt: string): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dn-pasek-godlo dn-powloka__godlo';

  const godlo = elementGodla({
    rozmiar: 24,
    podloze: 'ciemne',
    etykieta: 'Danaco Holding Group',
  });

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-powloka__nazwa';
  nazwa.textContent = NAZWA_PRODUKTU;

  element.append(godlo, nazwa, ...nazwaProjektu(projekt));
  return element;
}

/** Nazwa projektu jako wykaz — pusty, gdy powtarzałaby nazwę produktu. */
function nazwaProjektu(projekt: string): HTMLElement[] {
  const tresc = projekt.trim();
  if (tresc === '' || tresc.toLocaleLowerCase('pl') === NAZWA_PRODUKTU.toLocaleLowerCase('pl')) {
    return [];
  }
  const element = document.createElement('span');
  element.className = 'dn-powloka__projekt';
  element.textContent = tresc;
  return [element];
}
