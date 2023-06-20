package process

import (
	"fmt"
	"github.com/eryajf/chatgpt-dingtalk/pkg/db"
	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/solywsh/chatgpt"
	"reflect"
	"regexp"
	"strings"
)

// MjPrompt // Mj 工具
func MjPrompt(rmsg *dingbot.ReceiveMsg) error {
	qObj := db.Chat{
		Username:      rmsg.SenderNick,
		Source:        rmsg.GetChatTitle(),
		ChatType:      db.Q,
		ParentContent: 0,
		Content:       rmsg.Text.Content,
	}
	qid, err := qObj.Add()
	if err != nil {
		logger.Error("往MySQL新增数据失败,错误信息：", err)
	}

	// 1. 从消息 rmsg.text.content 中提取 版本号 v5 v4 niji testp
	msg := rmsg.Text.Content
	// 1.1 提取版本号
	title := strings.Split(msg, "#")[1]
	version := strings.Split(title, "@")[1]
	// 1.2 提取问题
	question := strings.Split(msg, "#")[2]
	// 去除换行符
	question = strings.ReplaceAll(question, "\n", "")
	logger.Info("version: ", version, "question: ", question)
	rmsg.Text.Content = question

	var reply string
	if version == "Testp" {
		reply, err = Testp(rmsg, nil, 35)
	} else if version == "v4" {
		reply, err = V4(rmsg, nil, 35)
	} else if version == "niji" {
		reply, err = Niji(rmsg, nil, 35)
	} else {
		reply, err = V5(rmsg, nil, 35)
	}

	if err != nil {
		logger.Info(fmt.Errorf("gpt request error: %v", err))
		if strings.Contains(fmt.Sprintf("%v", err), "maximum question length exceeded") {
			public.UserService.ClearUserSessionContext(rmsg.GetSenderIdentifier())
			_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v\n\n> 已超过最大文本限制，请缩短提问文字的字数。", err))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return err
			}
		} else {
			_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v", err))
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return err
			}
		}
	}
	if reply == "" {
		logger.Warning(fmt.Errorf("get gpt result falied: %v", err))
		return nil
	} else {
		reply = strings.TrimSpace(reply)
		reply = strings.Trim(reply, "\n")
		aObj := db.Chat{
			Username:      rmsg.SenderNick,
			Source:        rmsg.GetChatTitle(),
			ChatType:      db.A,
			ParentContent: qid,
			Content:       reply,
		}
		_, err := aObj.Add()
		if err != nil {
			logger.Error("往MySQL新增数据失败,错误信息：", err)
		}
		logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.SenderNick, reply))
		if public.JudgeSensitiveWord(reply) {
			reply = public.SolveSensitiveWord(reply)
		}
		// 回复@我的用户
		_, err = rmsg.ReplyToDingtalk(string(dingbot.MARKDOWN), FormatMarkdown(reply))
		if err != nil {
			logger.Warning(fmt.Errorf("send message error: %v", err))
			return err
		}
	}

	return nil
}

// Mj

var promptModels = map[string]string{
	"artistic": "Create a English NNN-word description on the topic: \"TEXT\".\nExamples: \n\"cyborg woman| with a visible detailed brain| muscles cable wires| detailed cyberpunk background with neon lights| biopunk| cybernetic| unreal engine| CGI | ultra detailed| 4k\".\n\"video games icons, 2d icons, rpg skills icons, world of warcraft items icons, league of legends items icons, ability icon, fantasy, potions, spells, objects, flowers, gems, swords, axe, hammer, fire, ice, arcane, shiny object, graphic design, high contrast, artstation - -uplight --v 4\",\n\"2 warrior princesses fight, Dynamic pose; Artgerm, Wlop, Greg Rutkowski; the perfect mix of Emily Ratajkowski, Ana de Armas, Kate Beckinsale, Kelly Brook and Adriana Lima as warrior princess; high detailed tanned skin; beautiful long hair, intricately detailed eyes; druidic leather vest; wielding an Axe; Attractive; Flames in background; Lumen Global Illumination, Lord of the Rings, Game of Thrones, Hyper-Realistic, Hyper-Detailed, 8k, --no watermarks --no cape --testp --ar 9:16 --v 4 --upbeta\",\n\"highly detailed matte painting stylized three quarters portrait of an anthropomorphic rugged happy fox with sunglasses! head animal person, background blur bokeh! ! ; by studio ghibli, makoto shinkai, by artgerm, by wlop, by greg rutkowski --v 4\",\n\"a photo of 8k ultra realistic archangel with 6 wings, full body, intricate purple and blue neon armor, ornate, cinematic lighting, trending on artstation, 4k, hyperrealistic, focused, high details, unreal engine 5, cinematic --ar 9:16 --s 1250 --q 2\",\n\"helldog, cerberus, red dog, big dog, 3 heads, rotten skin, background gate of hell, background hell, character design, side view, full body view, full body, aggressive stil, mythic still, unreal engine, hyper-detailed\",\n\"faeries in a bubble gum world , cinematic view, cinematic lighting, HD, 8k, unreal engine 5, full realistic, landscape, octane, high details, unreal engine 5, octane render, cryengine, volumetric lighting, cinematic, mood, dust, particles, particle effect, atmosphere, ray tracing, uhd, lighting rendered in Unreal Engine 5 and Cinema4k, Diffraction Grating, Ray Tracing Ambient Occlusion, Antialiasing, GLSL-Shaders, Post Process, Post Production, Cell Shading, Tone Mapping, photorealistic --v 4\",\n\"bronze metal raven, magical automaton, intricate details --v 4\",\n\"person = barbara palvin as undead necromancer lich. looks =dark straight hair, feminine figure, gorgeous, pretty face,beautiful body, wings, beautiful symmetrical glowing red eyes, two beautiful smooth silky legs:: pose = arms around body, Dragon with Massive Wings spread completely behind back:: shape = dance pose, curvy body, anatomically correct and fully showing symmetrical female chest, anatomically correct female torso and belly button, complete smooth pretty human arms and fingers, symmetrical face. clothing = revealing outfit, A pair of wearing a ivory breastplate and full-body-jewelry harness made from ice jewels :: environment = cemetary with frost aura ::details = highly detailed clothing,clear face that is fully shown, hyper quality style, top shelf jewelry with jewel arrays = dark fantasy, realistic, full female-body shot, both lower and, upper lips are completely down in clear detail, 3d printed:: artist = cgsociety, artgerm, trending on artstation, by victor titov. --video --s 3250 --seed 909101793 --ar 10:18 --q 1.5 --no amputees --no blur DOF --no fractals --v 3\",\n\"two lovers as neural networks embracing, beautiful, intricate details, cinematic lighting, beautiful concept art, surreal, art station\",\n\"realistic photograph of police nousr robot in modern city, cyberpunk, character design, detailed face, highly detailed, intricate details, symmetrical, smooth, digital 3d, hard surface, mechanical design, real-time, vfx, ultra hd, hdr\"\n\n\n\nAfter you create a description (not included in the word limit):\nIf in the description that you created a landscape described, then add \"landscape, many details, a lot of detail\" at the end.\nIf the description you created describes a portrait then add \"focused, blured background, body, portret\" to the end.\nIf the description you created describes a full-length person, then add \"full-body, legs, arms\" to the end.\nIf you are describing things that can have reflections, then add \"reflection\".\n\nAlways! Add at the end one of the types of lighting that suits better to the end  in format: \"Lighting: \". Here is the list: Spot Lighting, Rear light, Dark Light, Blinding Light, Candle Light, Concentrated Lighting, Twilight, Sunshine, Sunset, Lamp, Electric Arc, Moonlight, Neon Light, Night Light, Nuclear Light, Cinematic Light or similar.\nAlways! Add one of the styles that fits better to the end in format: \"Style: \". Here is the list: Fantasy, Dark Fantasy, Abstraction, Ethereal, Weirdcore, Dreampunk, Daydreampunk, Science, Surrealism, Unrealistic, Surreal, Realistic, Photorealism, Classic, Retro, Retrowave, Vintage, Cyberpunk, Punk, Modern, Futuristic, Sci-fi, Alchemy, Alien, Aurora, Magic, Mystic, Marvel Comics, Anime, Cartoon, Manga, Kawaii, Pastel, Neon, Aesthetic, Miniature.\nVery desirable. If some lighting suits the im.\nAfter the description, a more suitable style and type of lighting should always be written.\nAlways! Add 5 additional parameters in the format: \"Details: \". Here is a list of possible:: mood, dust, particles, particle effect, atmosphere, ray tracing, uhd, lighting, Diffraction Grating, Ray Tracing Ambient Occlusion, Antialiasing, GLSL-Shaders, Post Process, Post Production, Cell Shading, Tone Mapping and similar and any similar ones that will add the atmosphere to the description.\n\nAt the end it should be like this: \"A Victorian-style chair with chrome and ornate decorations reflects a distorted image in the water on the ground. The intricate patterns on the chair add an air of sophistication and elegance to any room. Lighting: Candle Light. Style: Victorian. Details: Cell Shading, atmosphere, ray tracing.\"\n",
	"weights":  "If you wanted to describe on English a picture in NNN words on this topic \"TEXT\", then what are the main details with their description.\nAnd evaluate each word according to the importance of their presence in the picture in the format \"word:: 1-10\", the total score can be from 1 to 10.\nDescribe in the format:\n\"the tree::10, star::6 clusters::3, a green::2 bird::7, anime::5, flowers::1\"\n\nAssign the most important details in your opinion to 10, the most unimportant ones to 1, and distribute the rest in this range also in order of importance.\nLeave only the details, do not send the description itself.\nIn response, do not add anything other than what I ask.\nDistribute them as much as possible in the range from 1 to 10.\nWrite everything in lower case.\n\nAnything further just add a comma:\nIf the theme means that something should not be or empty or without and etc. Mandatory \"--\" before \"no\", also the parameter cannot have weight \"::\". Only \"--no\" exists, there are no other options with \"--\".\nFor example:\n\"no hat::5\" -> \"--no hat\"\n\"no background\" -> \"--no background\"\n\"empty city\" -> \"--no people\"\n\"no trees\" -> \"--no trees\"\n\"without head\" -> \"--no head\"\n\"flowerless\" -> \"--no flower\"\nIt can't be like this \"no trees::5\", only like this \"--no trees\"\n\nIf in the description that you created a landscape described, then add \"landscape::5, many details::4, a lot of detail::5\" at the end.\n\nIf the description you created describes a portrait then add \"focused::3, blured background::2, body::3, portret::5\" to the end.\n\nIf the description you created describes a full-length person, then add \"full body::4, legs::4, arms::4\" to the end.\n\nIf something else, then nothing needs to be added.\n\nIf you are describing things that can have reflections, then add \"reflection::4\".\nIf you describe things that are related to space, then add \"space::10, star clusters::10\".\n\n\nVery desirable. If one style suits the image better than others, write \"name::5\". Example: \"cyberpunk::5\". Style Types: Fantasy, Dark Fantasy, Abstraction, Ethereal, Weirdcore, Dreampunk, Daydreampunk, Science, Surrealism, Unrealistic, Surreal, Realistic, Photorealism, Classic, Retro, Retrowave, Vintage, Cyberpunk, Punk, Modern, Futuristic, Sci-fi, Alchemy, Alien, Aurora, Magic, Mystic, Marvel Comics, Anime, Cartoon, Manga, Kawaii, Pastel, Neon, Aesthetic, Miniature.\nVery desirable. If some lighting suits the image better than others, write \"name::5\". Example: \"neon light::5\". Light Types: Spot Lighting, Rear light, Dark Light, Blinding Light, Candle Light, Concentrated Lighting, Twilight, Sunshine, Sunset, Lamp, Electric Arc, Moonlight, Neon Light, Night Light, Nuclear Light, Cinematic Light or similar\n\nIf something goes wrong. You don't need to report it, just ignore it.\nAt the end your answer should look something like this: \"plain meadow::4 clear sky::3 small rock::2 landscape::5 many details::4 --no flowers --no trees\"",
}

type Settings map[string]map[string]string

var urlRegex = regexp.MustCompile(`^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$`)
var noParamRegex = regexp.MustCompile(`(no\s\w+)\:\:\d+|--no\s\S+`)
var aspectRatioRegex = regexp.MustCompile(`^\d+:\d+$`)

func createPrompt(text string, model string, words int) (string, error) {
	prompt := promptModels[model]
	prompt = strings.ReplaceAll(prompt, "NNN", fmt.Sprintf("%d", words))
	prompt = strings.ReplaceAll(prompt, "TEXT", text)

	return prompt, nil
}

func sliceNoParameters(response string, noParamRegex *regexp.Regexp) (string, []string) {
	noParameters := []string{}
	found := noParamRegex.FindStringSubmatchIndex(response)
	for found != nil {
		start := found[0]
		end := found[1]

		if found[2] != -1 {
			noParameters = append(noParameters, fmt.Sprintf("--%s", response[found[2]:found[3]]))
		} else {
			noParameters = append(noParameters, response[start:end])
		}

		responseToRunes := []rune(response)
		for i := start; i < end; i++ {
			responseToRunes[i] = ' '
		}
		response = string(responseToRunes)
		found = noParamRegex.FindStringSubmatchIndex(response)
	}

	return response, noParameters
}

func checkParameterInConfigAndSettings(configName string, modelSettingsName string, config map[string]string, modelSettings Settings) bool {
	// 判断 configName 是否在 config 中
	// 判断 config[configName] 是否在 modelSettings[modelSettingsName] 中
	if _, ok := config[configName]; ok {
		if contains(modelSettings[modelSettingsName], config[configName]) {
			return true
		}
	}
	return false
}

func contains(slice any, item any) bool {
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)

		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == item {
				return true
			}
		}
	}

	return false
}

func addCustomParametersToPrompt(prompt string, config map[string]string, settings Settings, exclude []string) string {
	parameters := []map[string]string{
		{
			"config_name":   "renderer",
			"settings_name": "renderers",
		},
		{
			"config_name":   "content",
			"settings_name": "contents",
		},
		{
			"config_name":   "type",
			"settings_name": "types",
		},
		{
			"config_name":   "aspect_ratio",
			"settings_name": "aspect_ratios",
		},
	}

	for _, parameter := range parameters {
		if !contains(exclude, parameter["config_name"]) {
			if checkParameterInConfigAndSettings(parameter["config_name"], parameter["settings_name"], config, settings) {
				prompt += " " + settings[parameter["settings_name"]][config[parameter["config_name"]]]
			}
		}
	}

	return prompt
}

func checkAndAddColor(prompt string, config map[string]string) string {
	if color, ok := config["color"]; ok {
		prompt += " " + fmt.Sprintf("%s::10", color)
	}
	return prompt
}

func checkAndAddURL(prompt string, config map[string]string, urlRegex *regexp.Regexp) string {
	if url, ok := config["url"]; ok && urlRegex.MatchString(url) {
		prompt = url + " " + prompt
	}
	return prompt
}

func checkAndAddNoParameters(prompt string, noParameters []string) string {
	if len(noParameters) > 0 {
		prompt += " " + strings.Join(noParameters, " ")
	}
	return prompt
}

func V5(rmsg *dingbot.ReceiveMsg, config map[string]string, words int) (string, error) {
	if words == 0 {
		words = 35
	}
	setting := Settings{
		"prompt_models": {
			"weights":  "weights",
			"artistic": "artistic",
		},
		"types": {
			"anime":          "anime style::5 --upanime",
			"photorealistic": "high quality photo::5, soft light::2, sharp-focus::3, hyper realism::4",
			"avatar":         "high quality avatar::5, circle::5, square::5, sharp-focus::5",
			"couple avatar":  "anime::5, huge::4, kiss:4, high quality avatar::3, girl and boy::5, romantic::3",
		},
		"renderers": {
			"octane":        "octane render::4",
			"unreal engine": "unreal engine::4",
			"ray tracing":   "ray tracing::4, v-ray::4",
			"mixed":         "octane render::3 unreal engine::3, v-ray::3",
		},
		"contents": {
			"character": "character design::4",
			"landscape": "landscape design::4",
			"object":    "object design::4, one object::5, high object details::3, --no background-details, --no background",
			"light":     "light-design::4, shaders::3",
			"particles": "particle design::4, particles::3, smoke::2, neon::2, flash::2, spark::2",
		},
	}

	text := rmsg.Text.Content

	model := "weights"
	if checkParameterInConfigAndSettings("model", "prompt_models", config, setting) {
		model = setting["prompt_models"][config["model"]]
	}

	gptPrompt, _ := createPrompt(text, model, words)

	gptResponse, _ := chatgpt.MjPrompt(gptPrompt, rmsg.GetSenderIdentifier())

	mjPrompt, noParameters := sliceNoParameters(gptResponse, noParamRegex)

	mjPrompt = checkAndAddURL(mjPrompt, config, urlRegex)

	mjPrompt = checkAndAddColor(mjPrompt, config)

	mjPrompt = addCustomParametersToPrompt(mjPrompt, config, setting, []string{"aspect_ratio"})

	if aspectRatio, ok := config["aspect_ratio"]; ok && aspectRatioRegex.MatchString(aspectRatio) {
		mjPrompt += " --ar " + aspectRatio
	}

	mjPrompt += " --v 5 --s 1000 --q 2"

	mjPrompt = checkAndAddNoParameters(mjPrompt, noParameters)

	return mjPrompt, nil
}

func V4(rmsg *dingbot.ReceiveMsg, config map[string]string, words int) (string, error) {
	if words == 0 {
		words = 35
	}
	settings := Settings{
		"prompt_models": {
			"weights":  "weights",
			"artistic": "artistic",
		},
		"types": {
			"anime":          "anime style::5 --upanime",
			"photorealistic": "high quality photo::5, soft light::2, sharp-focus::3, hyper realism::4",
			"avatar":         "high quality avatar::5, circle::5, square::5, sharp-focus::5",
			"couple avatar":  "anime::5, huge::4, kiss:4, high quality avatar::3, girl and boy::5, romantic::3",
		},
		"contents": {
			"character": "character design::4",
			"landscape": "landscape design::4",
			"object":    "object design::4, one object::5, high object details::3, --no background-details, --no background",
			"light":     "light-design::4, shaders::3",
			"particles": "particle design::4, particles::3, smoke::2, neon::2, flash::2, spark::2",
		},
		"renderers": {
			"octane":        "octane render::4",
			"unreal engine": "unreal engine::4",
			"ray tracing":   "ray tracing::4, v-ray::4",
			"mixed":         "octane render::3 unreal engine::3, v-ray::3",
		},
		"aspect_ratios": {
			"1:1":  "--ar 1:1",
			"1:2":  "--ar 1:2",
			"2:1":  "--ar 2:1",
			"2:3":  "--ar 2:3",
			"3:2":  "--ar 3:2",
			"4:5":  "--ar 4:5",
			"5:4":  "--ar 5:4",
			"4:7":  "--ar 4:7",
			"7:4":  "--ar 7:4",
			"16:9": "--ar 16:9",
			"9:16": "--ar 9:16",
		},
	}
	model := "weights"
	if checkParameterInConfigAndSettings("model", "prompt_models", config, settings) {
		model = settings["prompt_models"][config["model"]]
	}

	text := rmsg.Text.Content

	gptPrompt, _ := createPrompt(text, model, words)

	gptResponse, _ := chatgpt.MjPrompt(gptPrompt, rmsg.GetSenderIdentifier())

	mjPrompt, noParameters := sliceNoParameters(gptResponse, noParamRegex)

	mjPrompt = checkAndAddURL(mjPrompt, config, urlRegex)

	mjPrompt = checkAndAddColor(mjPrompt, config)

	mjPrompt = addCustomParametersToPrompt(mjPrompt, config, settings, nil)

	mjPrompt += " --v 4 --s 1000 --q 5"

	mjPrompt = checkAndAddNoParameters(mjPrompt, noParameters)

	return mjPrompt, nil
}

func Niji(rmsg *dingbot.ReceiveMsg, config map[string]string, words int) (string, error) {
	if words == 0 {
		words = 35
	}
	settings := Settings{
		"prompt_models": {
			"weights":  "weights",
			"artistic": "artistic",
		},
		"contents": {
			"character": "character design::4",
			"landscape": "landscape design::4",
			"object":    "object design::4, one object::5, high object details::3, --no background-details, --no background",
			"light":     "light-design::4, shaders::3",
			"particles": "particle design::4, particles::3, smoke::2, neon::2, flash::2, spark::2",
		},
		"renderers": {
			"octane":        "octane render::4",
			"unreal engine": "unreal engine::4",
			"ray tracing":   "ray tracing::4, v-ray::4",
			"mixed":         "octane render::3 unreal engine::3, v-ray::3",
		},
		"aspect_ratios": {
			"1:2": "--ar 1:2",
			"2:1": "--ar 2:1",
		},
	}

	model := "weights"
	if checkParameterInConfigAndSettings("model", "prompt_models", config, settings) {
		model = settings["prompt_models"][config["model"]]
	}

	text := rmsg.Text.Content

	gptPrompt, _ := createPrompt(text, model, words)

	gptResponse, _ := chatgpt.MjPrompt(gptPrompt, rmsg.GetSenderIdentifier())

	mjPrompt, noParameters := sliceNoParameters(gptResponse, noParamRegex)

	mjPrompt = checkAndAddURL(mjPrompt, config, urlRegex)

	mjPrompt = checkAndAddColor(mjPrompt, config)

	mjPrompt = addCustomParametersToPrompt(mjPrompt, config, settings, []string{"type"})

	mjPrompt += " --niji --q 2"

	mjPrompt = checkAndAddNoParameters(mjPrompt, noParameters)

	return mjPrompt, nil
}

func Testp(rmsg *dingbot.ReceiveMsg, config map[string]string, words int) (string, error) {
	if words == 0 {
		words = 35
	}

	text := rmsg.Text.Content

	gptPrompt, _ := createPrompt(text, "artistic", words)

	gptResponse, _ := chatgpt.MjPrompt(gptPrompt, rmsg.GetSenderIdentifier())

	mjPrompt, noParameters := sliceNoParameters(gptResponse, noParamRegex)

	mjPrompt = checkAndAddURL(mjPrompt, config, urlRegex)

	settings := Settings{
		"aspect_ratios": {
			"2:3": "--ar 2:3",
			"3:2": "--ar 3:2",
		},
	}

	if checkParameterInConfigAndSettings("aspect_ratio", "aspect_ratios", config, settings) {
		mjPrompt += " " + settings["aspect_ratios"][config["aspect_ratio"]]
	}

	mjPrompt += " --testp --s 1500"

	mjPrompt = checkAndAddNoParameters(mjPrompt, noParameters)

	return mjPrompt, nil
}
