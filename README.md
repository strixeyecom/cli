[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/usestrix/cli">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">The StrixEye CLI</h3>

  <p align="center">
    Get the most out of your StrixEye experience
    <br />
    <a href="https://github.com/usestrix/cli"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/usestrix/cli">View Demo</a>
    ·
    <a href="https://github.com/usestrix/cli/issues">Report Bug</a>
    ·
    <a href="https://github.com/usestrix/cli/issues">Request Feature</a>
  </p>
</p>



<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgements">Acknowledgements</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

[![Product Name Screen Shot][product-screenshot]](https://strixeye.com)



### Built With

Thanks to maintainers and communities of the following projects for making development of our CLI easier. Full list of dependencies can be found in go modules file.
* [Cobra](https://github.com/spf13/cobra)
* [Viper](https://github.com/spf13/viper)
* [Laravel](https://laravel.com)



<!-- GETTING STARTED -->
## Getting Started

General information about setting up StrixEye CLI locally

### Prerequisites

Required softwares and installations.
* CLI has no external dependencies. It will work on all machines those operating systems we support.

### Installation

1. Register or contact [StrixEye](https://strixeye.com/)
1. Get your User API Key from [StrixEye Dashboard](https://dashboard.strixeye.com/settings/profile)
2. Get StrixEye CLI from your package manager or see <a href="#build">Build</a>
3. Authenticate yourself to StrixEye User API with CLI
   ```sh
   strixeye configure user
   ```
   or
   ```sh
   strixeye login
   ```

   Then, enter your User API Key when prompted.
   ![StrixEye CLI Login Processs](https://raw.githubusercontent.com/usestrix/cli/master/data/screenshots/login.gif)
   
   
4. Choose an agent to work with:
   ![StrixEye CLI Login Processs](https://raw.githubusercontent.com/usestrix/cli/master/data/screenshots/agents.gif)

   ```sh
   strixeye configure agent
   ```



<!-- USAGE EXAMPLES -->
## Usage



_For more examples, please refer to the [Documentation](https://example.com)_



<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/usestrix/cli/issues) for a list of proposed features (and known issues).



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the Apache License 2.0 License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

Your Name - [@strixeye](https://twitter.com/strixeye) - help@strixeye.com

Project Link: [https://github.com/usestrix/cli](https://github.com/usestrix/cli)




<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/usestrix/cli.svg?style=for-the-badge
[contributors-url]: https://github.com/usestrix/cli/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/usestrix/cli.svg?style=for-the-badge
[forks-url]: https://github.com/usestrix/clinetwork/members
[stars-shield]: https://img.shields.io/github/stars/usestrix/cli?style=for-the-badge
[stars-url]: https://github.com/usestrix/cli/stargazers
[issues-shield]: https://img.shields.io/github/issues/usestrix/cli.svg?style=for-the-badge
[issues-url]: https://github.com/usestrix/cli/issues
[license-shield]: https://img.shields.io/github/license/usestrix/cli.svg?style=for-the-badge
[license-url]: https://github.com/usestrix/cli/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/strixeye
[product-screenshot]: images/screenshot.png
