import '../motyw/motyw.css';
import '../komponenty/indeks.css';

import { ProgressStatus } from '../../../shared/contract';
import { uruchomMotyw, zastosujMotyw } from '../motyw/motyw';
import { utworzPowloke } from './powloka';

/**
 * Podgląd powłoki środowiska — strona sprawdzająca, wzorowana na
 * `motyw/podglad-zetonow`.
 *
 * Jedna odpowiedzialność: uruchomienie powłoki poza aplikacją, żeby cztery
 * pasy dało się obejrzeć i przeklikać przed ich osadzeniem w punkcie wejścia.
 * Podgląd niczego nie udaje: karty sesji i ich stany są tu materiałem
 * pokazowym, a nie danymi z rdzenia — dlatego mieszkają na tej stronie,
 * a nie w module powłoki.
 *
 * Uruchomienie:  npm run dev  →  /src/powloka/podglad.html
 * Motyw:         ?motyw=jasny  albo  ?motyw=ciemny; bez parametru rozstrzyga
 *                zapisany wybór, a w jego braku preferencja systemu.
 */

uruchomMotyw();

// Oba motywy są równoprawne, więc podgląd musi umieć pokazać każdy z nich
// bez ruszania ustawień środowiska.
const zadany = new URLSearchParams(window.location.search).get('motyw');
if (zadany === 'jasny') zastosujMotyw('light');
if (zadany === 'ciemny') zastosujMotyw('dark');

const powloka = utworzPowloke({ operator: 'Dariusz Naharnowicz' });

// Trzy dodatkowe karty pokazują wskaźnik pracy w tle w każdym z jego stanów.
powloka.karty.dodaj('Analiza rynku', ProgressStatus.Running);
powloka.karty.dodaj('Umowa ramowa', ProgressStatus.Paused);
powloka.karty.dodaj('Research OZE', ProgressStatus.Done);
powloka.pasek.ustawPowiadomienia(3);

const pierwsza = powloka.karty.wykaz()[0];
if (pierwsza !== undefined) powloka.karty.wybierz(pierwsza.id);

document.body.append(powloka.element);
