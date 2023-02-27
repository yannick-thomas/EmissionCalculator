import tkinter as tk
from tkinter import ttk
import customtkinter

from calculator.renderer.oil_renderer import OilRenderer
from calculator.renderer.briquettes_renderer import BriquettesRenderer


class CalcApp(customtkinter.CTk):

    WIDTH = 980
    HEIGHT = 720

    def __init__(self):
        super().__init__()

        self.delivery_quantity_briquettes = tk.StringVar(self)
        self.delivery_quantity_oil = tk.StringVar(self)
        self.result_briquettes = 0
        self.is_valid_calc = False

        # Window configuration
        self.title("Emissionsrechner")
        self.geometry(f"{CalcApp.WIDTH}x{CalcApp.HEIGHT}")
        self.minsize(self.WIDTH, self.HEIGHT)
        self.maxsize(self.WIDTH, self.HEIGHT)

        self.grid_columnconfigure(1, weight=1)
        self.grid_columnconfigure((2, 3), weight=0)
        self.grid_rowconfigure((0, 1, 2), weight=1)

        ## Sidebar
        self.sidebar = customtkinter.CTkFrame(self, width=160, corner_radius=0, fg_color="#404040")
        self.sidebar.grid(row=0, column=0, rowspan=5, sticky="nsew")
        self.sidebar.grid_rowconfigure(5, weight=1)
        self.logo_label = customtkinter.CTkLabel(self.sidebar, text="Emissionsrechner", font=customtkinter.CTkFont(size=25, weight="bold"))
        self.logo_label.grid(row=0, column=0, padx=20, pady=(20, 10))
        self.sidebar_briquettes = customtkinter.CTkButton(self.sidebar, command=self.calc_briquettes, text="Briketts")
        self.sidebar_briquettes.grid(row=1, column=0, padx=20, pady=10)
        self.sidebar_heating_oil = customtkinter.CTkButton(self.sidebar, command=self.calc_heating_oil, text="Heizöl")
        self.sidebar_heating_oil.grid(row=2, column= 0, padx= 20, pady=10)

    ## functions

    def calc_briquettes(self):
        briquettes_renderer = BriquettesRenderer(self)
        briquettes_renderer.build_base()
        
    def calc_heating_oil(self):
        oil_renderer = OilRenderer(self)
        oil_renderer.build_base()

if __name__ == "__main__":
    app = CalcApp()
    app.mainloop()