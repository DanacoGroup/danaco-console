# Pomocnik odtwarzania twarzy — jedyny kod tej rodziny liczący sieć twarzową.
#
# Rdzeń jest w Go, a GFPGAN wydany jest jako wagi PyTorcha (`GFPGANv1.4.pth`).
# Przepisanie tej sieci do Go byłoby drugą implementacją cudzej architektury
# i rozjeżdżałoby się z wagami przy każdym kolejnym wydaniu modelu, więc
# przebieg twarzowy jest procesem obok rdzenia — tak samo jak liczenie wektorów
# znaczenia (`internal/wiedza/pomocnik_osadzen.py`) i rozpoznawanie mowy
# (`internal/mowa/pomocnik.go`).
#
# Zlecenie przychodzi czterema argumentami wiersza poleceń, a odpowiedź wraca
# jednym obiektem JSON na standardowym wyjściu. Diagnostyka bibliotek idzie na
# strumień diagnostyczny, bo ostrzeżenie wstawione w środek JSON-a uczyniłoby
# odpowiedź nieczytelną. Standardowego wejścia nie ma — `zewnetrzne.Wolaj` go
# nie podaje.
#
# ── Dlaczego przebieg jest OSOBNY, a nie wpięty w powiększanie ───────────────
# Kontrakt `image.upscale` nazywa to polem `faces` opisanym jako „osobny
# przebieg" i tak też jest liczone: najpierw Real-ESRGAN powiększa cały obraz,
# potem ten pomocnik odnajduje w wyniku twarze, odtwarza każdą z osobna
# w rozdzielczości 512×512 i wkleja ją z powrotem. Wymiary wyniku pochodzą więc
# wyłącznie z powiększenia — sieć twarzowa ich nie rusza i odpowiedź kontraktu
# niesie te same `width` i `height`, co przebieg bez poprawki.
#
# ── Dlaczego GFPGAN, a nie CodeFormer ───────────────────────────────────────
# Obok `GFPGANv1.4.pth` leży `codeformer.pth`. Rdzeń go nie woła, bo CodeFormer
# stoi na własnej architekturze (VQGAN wraz z transformerem przewidującym kod
# słownika), której wydanie nie niesie w wagach — trzeba by wnieść drugi zestaw
# cudzego kodu obok tego, którym GFPGAN już liczy. Jedna sieć twarzowa, którą
# widać w treści odmowy rdzenia, jest tu wyborem świadomym, nie brakiem.
#
# ── Brak jest odpowiedzią, a nie wywróceniem ────────────────────────────────
# Gdy biblioteki nie ma albo wagi są nie do wczytania, pomocnik oddaje
# `{"ok": false, "powod": …}` i kończy pracę kodem zerowym; rdzeń zamienia to na
# odmowę nazywającą brak. Sam ślad stosu Pythona nie powiedziałby Operatorowi,
# czego brakuje.
import json
import sys

# WAGA_ODTWORZENIA steruje udziałem sieci w wyniku i jest wartością domyślną
# wydania GFPGAN (`gfpgan/utils.py`). Wyżej znaczy twarz gładszą i dalszą od
# oryginału, niżej — bliższą źródłu i słabiej poprawioną. Kontrakt nie ma pola
# na tę liczbę, więc rdzeń nie wystawia jej na zewnątrz i trzyma wartość autora
# sieci zamiast zgadywać własną.
WAGA_ODTWORZENIA = 0.5

# ROZMIAR_TWARZY to bok kwadratu, na którym pracuje sieć. Wynika z wag: GFPGAN
# v1.4 jest wyuczony na wycinkach 512×512 wyrównanych do pięciu punktów
# charakterystycznych i inna liczba nie pasuje do wymiarów warstw.
ROZMIAR_TWARZY = 512

# NAJMNIEJSZY_ROZSTAW_OCZU odrzuca znaleziska mniejsze niż pięć pikseli między
# źrenicami. Taki wycinek po rozciągnięciu do 512 pikseli jest samym szumem,
# a wklejony z powrotem zostawia w obrazie plamę wyraźniejszą od tego, co było.
NAJMNIEJSZY_ROZSTAW_OCZU = 5


def odpowiedz(tresc):
    """Wypisuje jeden obiekt JSON na standardowe wyjście i kończy pracę."""
    sys.stdout.write(json.dumps(tresc, ensure_ascii=False))
    sys.stdout.flush()
    sys.exit(0)


def main():
    if len(sys.argv) != 5:
        odpowiedz({"ok": False, "powod": "pomocnik przyjmuje cztery argumenty: "
                                         "wejście, wyjście, wagi GFPGAN, katalog wag pomocniczych"})
    wejscie, wyjscie, wagi, katalogWag = sys.argv[1:5]

    try:
        import cv2
        import torch
        from facexlib.utils.face_restoration_helper import FaceRestoreHelper
        from gfpgan_clean import GFPGANv1Clean
        from torchvision.transforms.functional import normalize
    except ImportError as blad:
        odpowiedz({"ok": False, "powod": "środowisko pomocnika nie ma biblioteki: " + str(blad)})

    obraz = cv2.imread(wejscie, cv2.IMREAD_COLOR)
    if obraz is None:
        odpowiedz({"ok": False, "powod": "nie da się odczytać obrazu wejściowego " + wejscie})

    # Liczymy na procesorze bezwarunkowo — tą samą drogą, co pozostałe silniki
    # tej rodziny. Wybór karty graficznej robiłby z jednego przebiegu dwa różne
    # w zależności od maszyny, a wynik ma być powtarzalny.
    urzadzenie = torch.device("cpu")

    try:
        siec = GFPGANv1Clean(
            out_size=ROZMIAR_TWARZY, num_style_feat=512, channel_multiplier=2,
            decoder_load_path=None, fix_decoder=False, num_mlp=8,
            input_is_latent=True, different_w=True, narrow=1, sft_half=True)
        zapis = torch.load(wagi, map_location="cpu", weights_only=True)
        siec.load_state_dict(zapis.get("params_ema", zapis), strict=True)
        siec.eval().to(urzadzenie)
    except Exception as blad:  # noqa: BLE001 — powód idzie do odmowy rdzenia
        odpowiedz({"ok": False, "powod": "wagi GFPGAN nie pasują do architektury: " + str(blad)})

    # `upscale_factor=1`, bo powiększenie zrobił już Real-ESRGAN i pomocnik
    # dostaje jego wynik. Każda inna wartość zmieniłaby wymiary obrazu po raz
    # drugi, a odpowiedź kontraktu obiecuje krotność podaną w żądaniu.
    pomocnicze = FaceRestoreHelper(
        1, face_size=ROZMIAR_TWARZY, crop_ratio=(1, 1),
        det_model="retinaface_resnet50", save_ext="png",
        use_parse=True, device=urzadzenie, model_rootpath=katalogWag)
    pomocnicze.read_image(obraz)
    znalezione = pomocnicze.get_face_landmarks_5(
        only_center_face=False, eye_dist_threshold=NAJMNIEJSZY_ROZSTAW_OCZU)
    pomocnicze.align_warp_face()

    for wycinek in pomocnicze.cropped_faces:
        tensor = torch.from_numpy(
            cv2.cvtColor(wycinek, cv2.COLOR_BGR2RGB).transpose(2, 0, 1).copy()).float() / 255.0
        normalize(tensor, (0.5, 0.5, 0.5), (0.5, 0.5, 0.5), inplace=True)
        with torch.no_grad():
            wynik = siec(tensor.unsqueeze(0).to(urzadzenie),
                         return_rgb=False, weight=WAGA_ODTWORZENIA)[0]
        wynik = wynik.squeeze(0).clamp(-1, 1).cpu().numpy().transpose(1, 2, 0)
        wynik = ((wynik + 1) / 2 * 255).round().clip(0, 255).astype("uint8")
        pomocnicze.add_restored_face(cv2.cvtColor(wynik, cv2.COLOR_RGB2BGR))

    pomocnicze.get_inverse_affine(None)
    zlozony = pomocnicze.paste_faces_to_input_image(upsample_img=None)

    # Plik powstaje ZAWSZE, także przy zerze znalezionych twarzy. Obraz bez
    # twarzy przechodzi wtedy nietknięty, a rdzeń oddaje wynik powiększenia —
    # odmowa byłaby tu karą za to, że na zdjęciu nikogo nie ma.
    if not cv2.imwrite(wyjscie, zlozony):
        odpowiedz({"ok": False, "powod": "nie da się zapisać wyniku pod " + wyjscie})
    odpowiedz({"ok": True, "twarze": int(znalezione)})


main()
