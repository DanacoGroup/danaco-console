import type { ConfigScope, WorkspaceInstructions } from '../../../../shared/contract';
import { nazwaPoziomu } from './poziomy-zasiegu';

/** Szablony instrukcji wstawiane do edytora, których treść stanowi materiał Operatora do dalszej redakcji. */
export const SZABLONY: ReadonlyArray<[string, string]> = [
  ['—', ''],
  [
    'Rola i zakres',
    '## Rola\nDziałasz w projekcie jako referent prowadzący.\n\n## Zakres\n- ...\n',
  ],
  [
    'Zasady redakcyjne',
    '## Zasady redakcyjne\n- Język polski, styl urzędowy.\n- Każde twierdzenie z podaniem podstawy.\n',
  ],
  [
    'Granice pracy',
    '## Granice\n- Nie wykonuj czynności nieodwracalnych bez decyzji Operatora.\n',
  ],
];


/**
 * Opis warstwy obowiązującej po zapisie: poziom, byt poziomu, odcisk treści
 * oraz ostrzeżenie, gdy zapisana warstwa została przykryta warstwą węższą.
 * Bez ostrzeżenia zapis wyglądałby na obowiązujący, choć obowiązuje treść
 * warstwy węższej.
 */
export function opisWarstwy(
  instrukcje: WorkspaceInstructions,
  zapisany: ConfigScope,
): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dw-warstwa';
  element.append(
    zdanie(`Warstwa obowiązująca: ${nazwaPoziomu(instrukcje.scope)}.`),
    zdanie(`Byt poziomu: ${instrukcje.scopeId ?? 'brak (poziom bez bytu)'}.`),
    zdanie(`Odcisk treści: ${instrukcje.contentHash ?? 'brak treści'}.`),
    zdanie(
      instrukcje.scope === zapisany
        ? 'Zapisana warstwa jest warstwą obowiązującą.'
        : 'UWAGA: zapisana warstwa jest przykryta warstwą węższą — poniżej treść obowiązująca.',
    ),
  );
  const obowiazujaca = document.createElement('pre');
  obowiazujaca.className = 'dw-warstwa__tresc';
  obowiazujaca.textContent = instrukcje.content;
  element.append(obowiazujaca);
  return element;
}

/** Buduje jedno zdanie opisu warstwy instrukcji, wyświetlane jako osobny akapit wewnątrz całego opisu warstwy. */
function zdanie(tekst: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = tekst;
  return element;
}
