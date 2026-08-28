# Pomocnik odtwarza twarze siecią GFPGAN jako proces obok rdzenia Go, ponieważ wagi
# PyTorcha nie dają się przepisać do Go bez rozjazdu z każdym wydaniem modelu.
import json
import sys

# WAGA_ODTWORZENIA steruje udziałem sieci w wyniku i jest wartością domyślną wydania
# GFPGAN. Wyżej znaczy twarz gładszą i dalszą od oryginału, niżej — bliższą źródłu
# i słabiej poprawioną.
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

    # Liczenie idzie na procesorze bezwarunkowo, dla powtarzalnego wyniku na każdej maszynie.
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

    # upscale_factor wynosi 1, bo powiększenie zrobił już Real-ESRGAN, a pomocnik pracuje na jego wyniku.
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

    # Plik wyniku powstaje zawsze, także bez znalezionych twarzy; obraz przechodzi wtedy nietknięty.
    if not cv2.imwrite(wyjscie, zlozony):
        odpowiedz({"ok": False, "powod": "nie da się zapisać wyniku pod " + wyjscie})
    odpowiedz({"ok": True, "twarze": int(znalezione)})


main()
