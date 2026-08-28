// Źródła SVG — grupa: srodowiska. Nazwy i kolejność wprost z `ikony/manifest.json`.
// Każdy znak wchodzi tu jako tekst dokumentu SVG i staje się wartością wykazu
// pod kluczem, którym reszta interfejsu ikonę przywołuje.

import srodowiskoTalkin from '../svg/srodowisko-talkin.svg?raw';
import srodowiskoWorkspace from '../svg/srodowisko-workspace.svg?raw';
import srodowiskoCodestudio from '../svg/srodowisko-codestudio.svg?raw';
import srodowiskoMultitaskingai from '../svg/srodowisko-multitaskingai.svg?raw';

/**
 * Emblematy czterech środowisk platformy — znaki własne, spoza zestawu Lucide.
 *
 * Wykaz jest zamrożony, więc zbiór nazw emblematów pozostaje zamknięty
 * i sprawdzian typów wychwytuje nazwę spoza niego.
 */
export const SRODOWISKA = {
  'srodowisko-talkin': srodowiskoTalkin,
  'srodowisko-workspace': srodowiskoWorkspace,
  'srodowisko-codestudio': srodowiskoCodestudio,
  'srodowisko-multitaskingai': srodowiskoMultitaskingai,
} as const;
